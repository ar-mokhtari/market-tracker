package usecase

import "github.com/ar-mokhtari/market-tracker/entity"

type Repo interface {
	Upsert(p entity.Price) error
	List(pType string) ([]entity.Price, error)
	GetHistory(symbol string, limit int) ([]entity.Price, error)
}
