package kt

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) DeleteKt(ctx context.Context, id uuid.UUID) error {
	return s.dao.NewKtQuery(ctx).DeleteKt(id)
}
