package validation

import (
	"errors"

	"github.com/ar-mokhtari/market-tracker/entity"
)

func ValidatePrice(p entity.Price) error {
	if p.Symbol == "" {
		return errors.New("symbol cannot be empty")
	}
	if p.Price == "" || p.Price == "0" {
		return errors.New("invalid price value")
	}
	return nil
}
