package mri

import (
	"context"

	"github.com/google/uuid"

	domain "composition-api/internal/domain/exam"
)

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (domain.Mri, error) {
	mri, err := s.adapters.Exam.GetMriById(ctx, id)
	if err != nil {
		return domain.Mri{}, err
	}
	return mri, nil
}

func (s *service) GetByExternalID(ctx context.Context, externalID uuid.UUID) ([]domain.Mri, error) {
	mris, err := s.adapters.Exam.GetMrisByExternalId(ctx, externalID)
	if err != nil {
		return nil, err
	}
	return mris, nil
}

func (s *service) GetByAuthor(ctx context.Context, author uuid.UUID) ([]domain.Mri, error) {
	mris, err := s.adapters.Exam.GetMrisByAuthor(ctx, author)
	if err != nil {
		return nil, err
	}
	return mris, nil
}

func (s *service) GetEchographicsByID(ctx context.Context, id uuid.UUID) (domain.Echographic, error) {
	echographics, err := s.adapters.Exam.GetEchographicByMriId(ctx, id)
	if err != nil {
		return domain.Echographic{}, err
	}
	return echographics, nil
}
