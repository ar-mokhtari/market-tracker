// Package mysql provides the database implementation of the repository.
package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ar-mokhtari/market-tracker/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Upsert updates the current price and adds to history only if the price has changed.
func (r *Repository) Upsert(p entity.Price) error {
	var lastPrice string

	// Check the latest price in history to avoid redundant records
	err := r.db.QueryRow("SELECT price FROM price_history WHERE symbol = ? ORDER BY recorded_at DESC LIMIT 1", p.Symbol).Scan(&lastPrice)

	// If price is different or it's the first record, proceed
	if err == sql.ErrNoRows || lastPrice != p.Price {
		tx, err := r.db.Begin()
		if err != nil {
			return fmt.Errorf("transaction start error: %w", err)
		}

		// 1. Update current prices table
		upsertQuery := `
			INSERT INTO prices (date, time, symbol, name_fa, price, unit, type)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
			price=VALUES(price), date=VALUES(date), time=VALUES(time), updated_at=CURRENT_TIMESTAMP`

		if _, err := tx.Exec(upsertQuery, p.Date, p.Time, p.Symbol, p.NameFa, p.Price, p.Unit, p.Type); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("upsert prices error: %w", err)
		}

		// 2. Insert into history
		if _, err := tx.Exec("INSERT INTO price_history (symbol, price, type) VALUES (?, ?, ?)", p.Symbol, p.Price, p.Type); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("insert history error: %w", err)
		}

		return tx.Commit()
	}

	return nil
}

func (r *Repository) List(pType string) ([]entity.Price, error) {
	query := `SELECT symbol, name_fa, price, unit, type, date, time
	          FROM prices WHERE type = ?`

	rows, err := r.db.Query(query, pType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []entity.Price
	for rows.Next() {
		var p entity.Price
		err := rows.Scan(&p.Symbol, &p.NameFa, &p.Price, &p.Unit, &p.Type, &p.Date, &p.Time)
		if err != nil {
			return nil, err
		}
		prices = append(prices, p)
	}
	return prices, nil
}

func (r *Repository) GetHistory(symbol string, limit int) ([]entity.Price, error) {
	query := `SELECT symbol, price, recorded_at
	          FROM price_history
	          WHERE symbol = ?
	          ORDER BY recorded_at DESC
	          LIMIT ?`

	rows, err := r.db.Query(query, symbol, limit)
	if err != nil {
		return nil, fmt.Errorf("repository history query error: %w", err)
	}
	defer rows.Close()

	var history []entity.Price
	for rows.Next() {
		var p entity.Price
		// Mapping recorded_at to Date for simplicity in this DTO
		var recordedAt time.Time
		if err := rows.Scan(&p.Symbol, &p.Price, &recordedAt); err != nil {
			return nil, fmt.Errorf("repository scan error: %w", err)
		}
		p.Date = recordedAt.Format("2006-01-02")
		p.Time = recordedAt.Format("15:04:05")
		history = append(history, p)
	}
	return history, nil
}
