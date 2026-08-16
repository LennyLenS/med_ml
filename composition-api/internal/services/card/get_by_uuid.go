package card

import (
	"context"

	domain "composition-api/internal/domain/med"

	"github.com/google/uuid"
)

func (s *service) GetCardByUUID(ctx context.Context, uuid uuid.UUID) (domain.Card, error) {
	card, err := s.adapters.Med.GetCardByUUID(ctx, uuid)
	if err != nil {
		return domain.Card{}, err
	}

	return card, nil
}
