package card

import (
	"context"
	"fmt"

	"med/internal/domain"

	"github.com/google/uuid"
)

func (s *service) GetCardByUUID(ctx context.Context, uuid uuid.UUID) (domain.Card, error) {
	card, err := s.dao.NewCardQuery(ctx).GetCardByUUID(uuid)
	if err != nil {
		return domain.Card{}, fmt.Errorf("get card by uuid: %w", err)
	}

	return card.ToDomain(), nil
}
