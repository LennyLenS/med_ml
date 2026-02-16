package card

import (
	"context"

	domain "composition-api/internal/domain/med"
)

func (s *service) GetCardByID(ctx context.Context, id int) (domain.Card, error) {
	card, err := s.adapters.Med.GetCardByID(ctx, id)
	if err != nil {
		return domain.Card{}, err
	}

	return card, nil
}
