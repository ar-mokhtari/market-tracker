package v1

import (
	"github.com/ar-mokhtari/market-tracker/usecase"
)

type Handler struct {
	UC *usecase.UseCase
}

func NewHandler(uc *usecase.UseCase) *Handler {
	return &Handler{UC: uc}
}
