package segment

import (
	"context"

	adapter "composition-api/internal/adapters/exam"
	domain "composition-api/internal/domain/exam"
)

func (s *service) Update(ctx context.Context, arg UpdateSegmentArg) (domain.Segment, error) {
	segment, err := s.adapters.Exam.UpdateSegment(ctx, adapter.UpdateSegmentIn{
		Id:        arg.Id,
		Contor:    arg.Contor,
		Knosp_012: arg.Knosp_012,
		Knosp_3:   arg.Knosp_3,
		Knosp_4:   arg.Knosp_4,
	})
	if err != nil {
		return domain.Segment{}, err
	}
	return segment, nil
}
