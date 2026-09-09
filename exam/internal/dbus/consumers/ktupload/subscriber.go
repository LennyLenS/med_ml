package ktupload

import (
	"context"
	"fmt"

	"github.com/WantBeASleep/med_ml_lib/dbus"
	"github.com/google/uuid"

	pb "exam/internal/generated/dbus/consume/ktupload"
	"exam/internal/services"
)

type subscriber struct {
	services *services.Services
}

func New(
	services *services.Services,
) dbus.Consumer[*pb.KtUpload] {
	return &subscriber{
		services: services,
	}
}

func (h *subscriber) Consume(ctx context.Context, event *pb.KtUpload) error {
	fmt.Println("kt upload event", event)
	if _, err := uuid.Parse(event.KtId); err != nil {
		return fmt.Errorf("kt id is not uuid: %s", event.KtId)
	}

	if err := h.services.Kt.PrepareKt(ctx, uuid.MustParse(event.KtId)); err != nil {
		return fmt.Errorf("process ktupload: %w", err)
	}

	return nil
}
