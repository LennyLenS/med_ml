package segment

import (
	"context"

	"github.com/google/uuid"

	adapter "composition-api/internal/adapters/exam"
)

func (s *service) Create(ctx context.Context, arg CreateSegmentArg) (uuid.UUID, error) {
	id, err := s.adapters.Exam.CreateSegment(ctx, adapter.CreateSegmentIn{
		ImageID:   arg.ImageID,
		NodeID:    arg.NodeID,
		Contor:    arg.Contor,
		Knosp_012: arg.Knosp_012,
		Knosp_3:   arg.Knosp_3,
		Knosp_4:   arg.Knosp_4,
	})
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
