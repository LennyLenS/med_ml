package image

import (
	"context"

	"github.com/google/uuid"

	domain "composition-api/internal/domain/exam"
)

func (s *service) GetImagesByMriID(ctx context.Context, mriID uuid.UUID) ([]domain.Image, error) {
	images, err := s.adapters.Exam.GetImagesByMriId(ctx, mriID)
	if err != nil {
		return nil, err
	}
	return images, nil
}
