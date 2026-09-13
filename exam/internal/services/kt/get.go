package kt

import (
	"context"

	"exam/internal/domain"

	ktEntity "exam/internal/repository/kt/entity"

	"github.com/google/uuid"
)

func (s *service) GetKtByID(ctx context.Context, id uuid.UUID) (domain.KT, error) {
	kt, err := s.dao.NewKtQuery(ctx).GetKtByID(id)
	if err != nil {
		return domain.KT{}, err
	}

	return kt.ToDomain(), nil
}

func (s *service) GetKtsByAuthor(ctx context.Context, author uuid.UUID) ([]domain.KT, error) {
	kts, err := s.dao.NewKtQuery(ctx).GetKtsByAuthor(author)
	if err != nil {
		return nil, err
	}

	return ktEntity.KT{}.SliceToDomain(kts), nil
}
