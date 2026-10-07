package kt

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	err := s.adapters.Exam.DeleteKt(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
