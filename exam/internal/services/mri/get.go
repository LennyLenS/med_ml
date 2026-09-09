package mri

import (
	"context"

	"exam/internal/domain"
	mriEntity "exam/internal/repository/mri/entity"

	"github.com/google/uuid"
)

func (s *service) GetMriByID(ctx context.Context, id uuid.UUID) (domain.Mri, error) {
	mri, err := s.dao.NewMriQuery(ctx).GetMriByID(id)
	if err != nil {
		return domain.Mri{}, err
	}

	return mri.ToDomain(), nil
}

func (s *service) GetMrisByExternalID(ctx context.Context, externalID uuid.UUID) ([]domain.Mri, error) {
	mris, err := s.dao.NewMriQuery(ctx).GetMrisByExternalID(externalID)
	if err != nil {
		return nil, err
	}

	return mriEntity.Mri{}.SliceToDomain(mris), nil
}

func (s *service) GetMrisByAuthor(ctx context.Context, author uuid.UUID) ([]domain.Mri, error) {
	mris, err := s.dao.NewMriQuery(ctx).GetMrisByAuthor(author)
	if err != nil {
		return nil, err
	}

	return mriEntity.Mri{}.SliceToDomain(mris), nil
}

func (s *service) GetMriEchographicsByID(ctx context.Context, id uuid.UUID) (domain.Echographic, error) {
	echographics, err := s.dao.NewEchographicQuery(ctx).GetEchographicByID(id)
	if err != nil {
		return domain.Echographic{}, err
	}

	return echographics.ToDomain(), nil
}
