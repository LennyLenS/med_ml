package kt

import (
	"context"
	"fmt"

	"exam/internal/domain"
	ktEntity "exam/internal/repository/kt/entity"
)

func (s *service) UpdateKt(ctx context.Context, arg UpdateKTArg) (domain.KT, error) {
	kt, err := s.GetKtByID(ctx, arg.Id)
	if err != nil {
		return domain.KT{}, fmt.Errorf("get kt by id: %w", err)
	}
	arg.UpdateDomain(&kt)

	if err := s.dao.NewKtQuery(ctx).UpdateKt(ktEntity.KT{}.FromDomain(kt)); err != nil {
		return domain.KT{}, fmt.Errorf("update kt: %w", err)
	}

	return kt, nil
}
