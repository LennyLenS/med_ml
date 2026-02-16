package card

import (
	"context"
	"fmt"

	"med/internal/domain"
)

func (s *service) GetCardByID(ctx context.Context, id int) (domain.Card, error) {
	card, err := s.dao.NewCardQuery(ctx).GetCardByID(id)
	if err != nil {
		return domain.Card{}, fmt.Errorf("get card by id: %w", err)
	}

	return card.ToDomain(), nil
}
