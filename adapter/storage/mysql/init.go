// Package mysql implements person storage with MySQL database.
package mysql

import (
	"database/sql"
)

type RawMySQLRepo struct {
	db *sql.DB
}

func NewRawRepo(db *sql.DB) *RawMySQLRepo {
	return &RawMySQLRepo{db: db}
}
