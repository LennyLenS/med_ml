package image

import (
	"context"

	"github.com/google/uuid"

	"composition-api/internal/adapters"
	domain "composition-api/internal/domain/exam"
)

type Service interface {
	GetImagesByMriID(ctx context.Context, mriID uuid.UUID) ([]domain.Image, error)
}

type service struct {
	adapters *adapters.Adapters
}

func New(
	adapters *adapters.Adapters,
) Service {
	return &service{
		adapters: adapters,
	}
}
