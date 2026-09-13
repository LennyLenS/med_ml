package device

import (
	"context"

	"exam/internal/domain"
	"exam/internal/repository"
)

type Service interface {
	CreateDevice(ctx context.Context, deviceName string) (int, error)

	GetDeviceList(ctx context.Context) ([]domain.Device, error)
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
