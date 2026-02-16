package card

import (
	"context"
	"errors"
	"log/slog"

	pb "med/internal/generated/grpc/service"
	"med/internal/repository/entity"
	"med/internal/server/mappers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) GetCardByID(ctx context.Context, in *pb.GetCardByIDIn) (*pb.GetCardOut, error) {
	slog.Info("GetCardByID called", "id", in.Id)

	if in.Id <= 0 {
		slog.Error("Invalid card ID", "id", in.Id)
		return nil, status.Errorf(codes.InvalidArgument, "Неверный ID карточки: должен быть положительным числом")
	}

	slog.Info("Calling service GetCardByID", "id", in.Id)
	card, err := h.cardSrv.GetCardByID(ctx, int(in.Id))
	if err != nil {
		slog.Error("Error getting card by ID", "id", in.Id, "err", err)
		switch {
		case errors.Is(err, entity.ErrNotFound):
			slog.Warn("Card not found", "id", in.Id)
			return nil, status.Errorf(codes.NotFound, "Карта не найдена")
		default:
			return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
		}
	}

	slog.Info("Card retrieved successfully", "id", in.Id, "doctorId", card.DoctorID, "patientId", card.PatientID, "diagnosis", card.Diagnosis)
	return &pb.GetCardOut{Card: mappers.CardFromDomain(card)}, nil
}
