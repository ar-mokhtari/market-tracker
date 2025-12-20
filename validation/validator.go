package validation

import (
	"errors"
	"github.com/ar-mokhtari/market-tracker/entity"
)

func ValidatePrice(p entity.Price) error {
	if p.Symbol == "" {
		return errors.New("نماد قیمت نمی‌تواند خالی باشد")
	}
	if p.Price == "" || p.Price == "0" {
		return errors.New("مقدار قیمت نامعتبر است")
	}
	return nil
}
