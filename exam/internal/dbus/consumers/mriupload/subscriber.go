package mriupload

import (
	"context"
	"fmt"

	"github.com/WantBeASleep/med_ml_lib/dbus"
	"github.com/google/uuid"

	pb "exam/internal/generated/dbus/consume/mriupload"
	"exam/internal/services"
)

type subscriber struct {
	services *services.Services
}

func New(
	services *services.Services,
) dbus.Consumer[*pb.MriUpload] {
	return &subscriber{
		services: services,
	}
}

func (h *subscriber) Consume(ctx context.Context, event *pb.MriUpload) error {
	if _, err := uuid.Parse(event.MriId); err != nil {
		return fmt.Errorf("mri id is not uuid: %s", event.MriId)
	}

	if err := h.services.Image.SplitMri(ctx, uuid.MustParse(event.MriId)); err != nil {
		return fmt.Errorf("process mriupload: %w", err)
	}
	return nil
}
