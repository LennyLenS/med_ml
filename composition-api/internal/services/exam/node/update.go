package node

import (
	"context"

	adapter "composition-api/internal/adapters/exam"
	domain "composition-api/internal/domain/exam"
)

func (s *service) UpdateNode(ctx context.Context, arg UpdateNodeArg) (domain.Node, error) {
	node, err := s.adapters.Exam.UpdateNode(ctx, adapter.UpdateNodeIn{
		Id:         arg.Id,
		Validation: arg.Validation,
		Knosp_012:  arg.Knosp_012,
		Knosp_3:    arg.Knosp_3,
		Knosp_4:    arg.Knosp_4,
	})
	if err != nil {
		return domain.Node{}, err
	}
	return node, nil
}
