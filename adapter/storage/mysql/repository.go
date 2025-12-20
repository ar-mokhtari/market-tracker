package mysql

import (
	"database/sql"

	"github.com/ar-mokhtari/market-tracker/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Upsert(p entity.Price) error {
	query := `INSERT INTO prices (date, time, symbol, name_en, name_fa, price, change_value, change_percent, unit, type)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			  ON DUPLICATE KEY UPDATE
			  price=VALUES(price), change_value=VALUES(change_value), change_percent=VALUES(change_percent),
			  date=VALUES(date), time=VALUES(time), updated_at=CURRENT_TIMESTAMP`

	_, err := r.db.Exec(query, p.Date, p.Time, p.Symbol, p.NameEn, p.NameFa, p.Price, p.ChangeValue, p.ChangePercent, p.Unit, p.Type)
	return err
}

func (r *Repository) List(pType string) ([]entity.Price, error) {
	query := "SELECT id, symbol, name_fa, price, unit, type, updated_at FROM prices"
	var rows *sql.Rows
	var err error

	if pType != "" {
		query += " WHERE type = ?"
		rows, err = r.db.Query(query, pType)
	} else {
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []entity.Price
	for rows.Next() {
		var p entity.Price
		rows.Scan(&p.ID, &p.Symbol, &p.NameFa, &p.Price, &p.Unit, &p.Type, &p.UpdatedAt)
		prices = append(prices, p)
	}
	return prices, nil
}

func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM prices WHERE id = ?", id)
	return err
}

func (r *Repository) GetByID(id string) (entity.Price, error) {
	var p entity.Price
	query := "SELECT id, symbol, name_fa, price, type FROM prices WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Symbol, &p.NameFa, &p.Price, &p.Type)
	return p, err
}

func (r *Repository) UpdatePrice(id string, newPrice string) error {
	_, err := r.db.Exec("UPDATE prices SET price = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", newPrice, id)
	return err
}
