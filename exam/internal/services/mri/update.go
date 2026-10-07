package mri

import (
	"context"
	"fmt"

	"exam/internal/domain"
	echographicEntity "exam/internal/repository/echographic/entity"
	mriEntity "exam/internal/repository/mri/entity"
)

func (s *service) UpdateMri(ctx context.Context, arg UpdateMriArg) (domain.Mri, error) {
	mri, err := s.GetMriByID(ctx, arg.Id)
	if err != nil {
		return domain.Mri{}, fmt.Errorf("get mri by id: %w", err)
	}
	arg.UpdateDomain(&mri)

	if err := s.dao.NewMriQuery(ctx).UpdateMri(mriEntity.Mri{}.FromDomain(mri)); err != nil {
		return domain.Mri{}, fmt.Errorf("update mri: %w", err)
	}

	return mri, nil
}

func (s *service) UpdateEchographic(ctx context.Context, arg UpdateEchographicArg) (domain.Echographic, error) {
	echographic, err := s.GetMriEchographicsByID(ctx, arg.Id)
	if err != nil {
		return domain.Echographic{}, fmt.Errorf("get mri by id: %w", err)
	}
	arg.UpdateDomain(&echographic)

	if err := s.dao.NewEchographicQuery(ctx).UpdateEchographic(echographicEntity.Echographic{}.FromDomain(echographic)); err != nil {
		return domain.Echographic{}, fmt.Errorf("update echographic: %w", err)
	}

	return echographic, nil
}
