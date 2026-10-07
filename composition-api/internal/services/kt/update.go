package kt

import (
	"context"

	adapter "composition-api/internal/adapters/exam"
	domain "composition-api/internal/domain/exam"
)

func (s *service) Update(ctx context.Context, arg UpdateKtArg) (domain.KT, error) {
	kt, err := s.adapters.Exam.UpdateKt(ctx, adapter.UpdateKtIn{
		Id:                 arg.Id,
		Checked:            arg.Checked,
		ClassProbabilities: arg.ClassProbabilities,
	})
	if err != nil {
		return domain.KT{}, err
	}
	return kt, nil
}
