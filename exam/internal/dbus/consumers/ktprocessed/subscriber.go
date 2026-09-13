package ktprocessed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/WantBeASleep/med_ml_lib/dbus"
	"github.com/google/uuid"

	pb "exam/internal/generated/dbus/consume/ktprocessed"
	"exam/internal/services"
)

var ErrInvalidContor = errors.New("invalid contor")

type subscriber struct {
	services *services.Services
}

func New(
	services *services.Services,
) dbus.Consumer[*pb.KtProcessed] {
	return &subscriber{
		services: services,
	}
}

func (h *subscriber) Consume(ctx context.Context, message *pb.KtProcessed) error {
	if _, err := uuid.Parse(message.KtId); err != nil {
		return fmt.Errorf("kt id is not uuid: %s", message.KtId)
	}

	probsJSON, err := json.Marshal(message.ClassProbabilities)
	if err != nil {
		return fmt.Errorf("failed to serialize: %v", err)
	}

	return h.services.Kt.SaveProcessedKt(ctx, uuid.MustParse(message.KtId), probsJSON)
}
