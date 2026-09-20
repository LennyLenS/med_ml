package mri

import (
	"context"

	adapter "composition-api/internal/adapters/exam"
	domain "composition-api/internal/domain/exam"
)

func (s *service) Update(ctx context.Context, arg UpdateMriArg) (domain.Mri, error) {
	mri, err := s.adapters.Exam.UpdateMri(ctx, adapter.UpdateMriIn{
		Id:         arg.Id,
		Projection: arg.Projection,
		Checked:    arg.Checked,
	})
	if err != nil {
		return domain.Mri{}, err
	}
	return mri, nil
}

func (s *service) UpdateEchographics(ctx context.Context, arg domain.Echographic) (domain.Echographic, error) {
	echographics, err := s.adapters.Exam.UpdateEchographic(ctx, arg)
	if err != nil {
		return domain.Echographic{}, err
	}
	return echographics, nil
}
