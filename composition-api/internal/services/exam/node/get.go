package node

import (
	"context"

	"github.com/google/uuid"

	domain "composition-api/internal/domain/exam"
)

func (s *service) GetNodesByMriID(ctx context.Context, mriID uuid.UUID) ([]domain.Node, error) {
	nodes, err := s.adapters.Exam.GetNodesByMriId(ctx, mriID)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
