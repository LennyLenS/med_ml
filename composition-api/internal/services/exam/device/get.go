package device

import (
	"context"

	domain "composition-api/internal/domain/exam"
)

func (s *service) GetAll(ctx context.Context) ([]domain.Device, error) {
	devices, err := s.adapters.Exam.GetDeviceList(ctx)
	if err != nil {
		return nil, err
	}
	return devices, nil
}
