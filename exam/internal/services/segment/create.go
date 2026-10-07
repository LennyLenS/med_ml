package segment

import (
	"context"
	"fmt"

	"exam/internal/domain"
	segmentEntity "exam/internal/repository/segment/entity"

	"github.com/google/uuid"
)

func (s *service) CreateManualSegment(ctx context.Context, arg CreateSegmentArg) (uuid.UUID, error) {
	node, err := s.dao.NewNodeQuery(ctx).GetNodeByID(arg.NodeID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("get node by id: %w", err)
	}

	if node.Ai {
		return uuid.Nil, ErrAddSegmentToAiNode
	}

	segment := domain.Segment{
		Id:       uuid.New(),
		ImageID:  arg.ImageID,
		NodeID:   arg.NodeID,
		Contor:   arg.Contor,
		Ai:       false,
		Knosp012: arg.Knosp012,
		Knosp3:   arg.Knosp3,
		Knosp4:   arg.Knosp4,
	}
	if err := s.dao.NewSegmentQuery(ctx).InsertSegments(segmentEntity.Segment{}.FromDomain(segment)); err != nil {
		return uuid.Nil, err
	}

	return segment.Id, nil
}
