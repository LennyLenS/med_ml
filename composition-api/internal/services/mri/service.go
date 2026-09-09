package mri

import (
	"context"

	"github.com/google/uuid"

	"composition-api/internal/adapters"
	dbus "composition-api/internal/dbus/producers"
	domain "composition-api/internal/domain/exam"
	"composition-api/internal/repository"
)

type Service interface {
	Create(ctx context.Context, arg CreateMriArg) (uuid.UUID, error)

	GetByID(ctx context.Context, id uuid.UUID) (domain.Mri, error)
	GetByExternalID(ctx context.Context, externalID uuid.UUID) ([]domain.Mri, error)
	GetByAuthor(ctx context.Context, author uuid.UUID) ([]domain.Mri, error)
	GetEchographicsByID(ctx context.Context, id uuid.UUID) (domain.Echographic, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error

	Update(ctx context.Context, arg UpdateMriArg) (domain.Mri, error)
	UpdateEchographics(ctx context.Context, arg domain.Echographic) (domain.Echographic, error)
}

type service struct {
	adapters *adapters.Adapters
	dao      repository.DAO
	dbus     dbus.Producer
}

func New(
	adapters *adapters.Adapters,
	dao repository.DAO,
	dbus dbus.Producer,
) Service {
	return &service{
		adapters: adapters,
		dao:      dao,
		dbus:     dbus,
	}
}
