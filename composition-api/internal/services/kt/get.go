package kt

import (
	"context"

	domain "composition-api/internal/domain/exam"

	"github.com/google/uuid"
)

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (domain.KT, error) {
	return domain.KT{}, nil
}
