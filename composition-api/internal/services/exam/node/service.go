package node

import (
	"context"

	"github.com/google/uuid"

	"composition-api/internal/adapters"
	domain "composition-api/internal/domain/exam"
)

type Service interface {
	GetNodesByMriID(ctx context.Context, mriID uuid.UUID) ([]domain.Node, error)
	UpdateNode(ctx context.Context, arg UpdateNodeArg) (domain.Node, error)
	DeleteNode(ctx context.Context, id uuid.UUID) error
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
