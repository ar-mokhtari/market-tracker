// Package mysql provides the database implementation of the repository.
package mysql

import (
	"context"
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

// Upsert ensures data is only added to history if the price has changed.
func (r *Repository) Upsert(p entity.Price) error {
	var lastPrice string

	// Query the most recent price for this specific symbol
	err := r.db.QueryRow("SELECT price FROM price_history WHERE symbol = ? ORDER BY recorded_at DESC LIMIT 1", p.Symbol).Scan(&lastPrice)

	// If price is identical, skip history insert but update the main prices table
	if err == nil && lastPrice == p.Price.String() {
		_, err = r.db.Exec("UPDATE prices SET date=?, time=?, updated_at=CURRENT_TIMESTAMP WHERE symbol=?", p.Date, p.Time, p.Symbol)
		return err
	}

	// Otherwise, record the change in both tables
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, _ = tx.Exec(`INSERT INTO prices (symbol, name_fa, price, unit, type, date, time)
		VALUES (?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE price=VALUES(price), date=VALUES(date), time=VALUES(time)`,
		p.Symbol, p.NameFa, p.Price, p.Unit, p.Type, p.Date, p.Time)

	_, _ = tx.Exec("INSERT INTO price_history (symbol, price, type) VALUES (?, ?, ?)", p.Symbol, p.Price, p.Type)

	return tx.Commit()
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

func (r *Repository) GetAllPrices(ctx context.Context, priceType string) ([]entity.Price, error) {
	var prices []entity.Price
	var query string
	var args []interface{}

	if priceType != "" {
		query = "SELECT date, time, symbol, price, type FROM prices WHERE type = ? ORDER BY created_at DESC"
		args = append(args, priceType)
	} else {
		query = "SELECT date, time, symbol, price, type FROM prices ORDER BY created_at DESC"
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p entity.Price
		// Make sure these fields match your entity.Price struct fields
		err := rows.Scan(&p.Date, &p.Time, &p.Symbol, &p.Price, &p.Type)
		if err != nil {
			return nil, err
		}
		prices = append(prices, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}
