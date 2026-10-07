package node_segment

import (
	"context"

	"github.com/google/uuid"

	domain "composition-api/internal/domain/exam"
)

func (s *service) GetNodeWithSegmentsByImageID(ctx context.Context, imageID uuid.UUID) ([]domain.Node, []domain.Segment, error) {
	nodes, segments, err := s.adapters.Exam.GetNodesWithSegmentsByImageId(ctx, imageID)
	if err != nil {
		return nil, nil, err
	}
	return nodes, segments, nil
}
