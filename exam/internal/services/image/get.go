package image

import (
	"context"

	"exam/internal/domain"
	"exam/internal/repository/image/entity"

	"github.com/google/uuid"
)

func (s *service) GetImagesByMriID(ctx context.Context, id uuid.UUID) ([]domain.Image, error) {
	images, err := s.dao.NewImageQuery(ctx).GetImagesByMriID(id)
	if err != nil {
		return nil, err
	}

	return entity.Image{}.SliceToDomain(images), nil
}
