// Package usecase is about manage and control services
package usecase

import (
	"github.com/ar-mokhtari/market-tracker/adapter/storage/mysql"
)

type UseCase struct {
	repo PriceRepository
}

func New(repo PriceRepository) *UseCase {
	return &UseCase{repo: repo}
}

type MySQLRepoAdapter struct {
	raw *mysql.RawMySQLRepo
}

func NewRepoAdapter(raw *mysql.RawMySQLRepo) *MySQLRepoAdapter {
	return &MySQLRepoAdapter{raw: raw}
}
