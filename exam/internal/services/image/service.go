package image

import (
	"context"

	dbus "exam/internal/dbus/producers"

	"exam/internal/domain"
	"exam/internal/repository"

	"github.com/google/uuid"
)

type Service interface {
	SplitMri(ctx context.Context, id uuid.UUID) error

	GetImagesByMriID(ctx context.Context, id uuid.UUID) ([]domain.Image, error)
}

type service struct {
	dao  repository.DAO
	dbus dbus.Producer
}

func New(
	dao repository.DAO,
	dbus dbus.Producer,
) Service {
	return &service{
		dao:  dao,
		dbus: dbus,
	}
}
