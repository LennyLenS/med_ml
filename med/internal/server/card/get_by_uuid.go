package card

import (
	"context"
	"errors"
	"log/slog"

	pb "med/internal/generated/grpc/service"
	"med/internal/repository/entity"
	"med/internal/server/mappers"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) GetCardByUUID(ctx context.Context, in *pb.GetCardByUUIDIn) (*pb.GetCardOut, error) {
	slog.Info("GetCardByUUID called", "uuid", in.Uuid)

	cardUUID, err := uuid.Parse(in.Uuid)
	if err != nil {
		slog.Error("Invalid card UUID", "uuid", in.Uuid, "err", err)
		return nil, status.Errorf(codes.InvalidArgument, "Неверный UUID карточки: %s", err.Error())
	}

	slog.Info("Calling service GetCardByUUID", "uuid", cardUUID)
	card, err := h.cardSrv.GetCardByUUID(ctx, cardUUID)
	if err != nil {
		slog.Error("Error getting card by UUID", "uuid", cardUUID, "err", err)
		switch {
		case errors.Is(err, entity.ErrNotFound):
			slog.Warn("Card not found", "uuid", cardUUID)
			return nil, status.Errorf(codes.NotFound, "Карта не найдена")
		default:
			return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
		}
	}

	slog.Info("Card retrieved successfully", "uuid", cardUUID, "doctorId", card.DoctorID, "patientId", card.PatientID)
	return &pb.GetCardOut{Card: mappers.CardFromDomain(card)}, nil
}
