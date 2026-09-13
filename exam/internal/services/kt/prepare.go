package kt

import (
	"context"
	"fmt"

	"exam/internal/generated/dbus/produce/ktprepared"

	"github.com/google/uuid"
)

func (s *service) PrepareKt(ctx context.Context, id uuid.UUID) error {
	if err := s.dbus.SendKtPrepared(ctx, &ktprepared.KtPrepared{
		KtId: id.String(),
	}); err != nil {
		return fmt.Errorf("send to ktprepared topic: %w", err)
	}

	return nil
}
