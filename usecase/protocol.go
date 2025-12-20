package usecase

import (
	"context"

	"github.com/ar-mokhtari/market-tracker/entity"
)

type PriceRepository interface {
	Create(ctx context.Context, p entity.Price) (int, error)
	GetByID(ctx context.Context, id int) (entity.Price, error)
	GetAll(ctx context.Context) ([]entity.Price, error)
	Update(ctx context.Context, p entity.Price) error
	Delete(ctx context.Context, id int) error
}
