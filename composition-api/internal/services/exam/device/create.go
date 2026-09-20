package device

import (
	"context"
)

func (s *service) Create(ctx context.Context, name string) (int, error) {
	id, err := s.adapters.Exam.CreateDevice(ctx, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}
