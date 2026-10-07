package mri

import (
	"context"

	"exam/internal/domain"
	"exam/internal/repository"

	"github.com/google/uuid"
)

type Service interface {
	CreateMri(ctx context.Context, arg CreateMriArg) (uuid.UUID, error)

	GetMriByID(ctx context.Context, id uuid.UUID) (domain.Mri, error)
	GetMrisByExternalID(ctx context.Context, externalID uuid.UUID) ([]domain.Mri, error)
	GetMrisByAuthor(ctx context.Context, author uuid.UUID) ([]domain.Mri, error)
	GetMriEchographicsByID(ctx context.Context, id uuid.UUID) (domain.Echographic, error)

	UpdateMri(ctx context.Context, arg UpdateMriArg) (domain.Mri, error)
	UpdateEchographic(ctx context.Context, arg UpdateEchographicArg) (domain.Echographic, error)

	DeleteMri(ctx context.Context, id uuid.UUID) error
}

type service struct {
	dao repository.DAO
}

func New(
	dao repository.DAO,
) Service {
	return &service{
		dao: dao,
	}
}
