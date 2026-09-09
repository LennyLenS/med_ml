package device

import (
	"context"

	"exam/internal/domain"
	"exam/internal/repository/device/entity"
)

func (s *service) GetDeviceList(ctx context.Context) ([]domain.Device, error) {
	devices, err := s.dao.NewDeviceQuery(ctx).GetDeviceList()
	if err != nil {
		return nil, err
	}

	return entity.Device{}.SliceToDomain(devices), nil
}
