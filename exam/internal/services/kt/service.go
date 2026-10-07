package kt

import (
	"context"

	dbus "exam/internal/dbus/producers"
	"exam/internal/domain"
	"exam/internal/repository"

	"github.com/google/uuid"
)

type Service interface {
	GetKtByID(ctx context.Context, id uuid.UUID) (domain.KT, error)
	GetKtsByAuthor(ctx context.Context, author uuid.UUID) ([]domain.KT, error)
	CreateKt(ctx context.Context, arg CreateKTArg) (uuid.UUID, error)
	PrepareKt(ctx context.Context, id uuid.UUID) error
	SaveProcessedKt(ctx context.Context, ktID uuid.UUID, arg []byte) error
	UpdateKt(ctx context.Context, arg UpdateKTArg) (domain.KT, error)
	DeleteKt(ctx context.Context, id uuid.UUID) error
}

type service struct {
	dao  repository.DAO
	dbus dbus.Producer
}

func New(
	dao repository.DAO,
	dbus dbus.Producer,
) Service {
	return &service{
		dao:  dao,
		dbus: dbus,
	}
}
