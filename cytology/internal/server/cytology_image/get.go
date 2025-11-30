package cytology_image

import (
	"context"
	"errors"

	"cytology/internal/domain"
	pb "cytology/internal/generated/grpc/service"
	"cytology/internal/server/mappers"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) GetCytologyImageById(ctx context.Context, in *pb.GetCytologyImageByIdIn) (*pb.GetCytologyImageByIdOut, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id is not a valid uuid: %s", err.Error())
	}

	img, err := h.services.CytologyImage.GetCytologyImageByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "cytology image not found")
		}
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	// Получаем оригинальное изображение если есть
	originalImages, err := h.services.OriginalImage.GetOriginalImagesByCytologyID(ctx, id)
	var originalImage *pb.OriginalImage
	if err == nil && len(originalImages) > 0 {
		originalImage = mappers.OriginalImageToProto(originalImages[0])
	}

	return &pb.GetCytologyImageByIdOut{
		CytologyImage: mappers.CytologyImageToProto(img),
		OriginalImage: originalImage,
	}, nil
}

func (h *handler) GetCytologyImagesByExternalId(ctx context.Context, in *pb.GetCytologyImagesByExternalIdIn) (*pb.GetCytologyImagesByExternalIdOut, error) {
	externalID, err := uuid.Parse(in.ExternalId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "external_id is not a valid uuid: %s", err.Error())
	}

	images, err := h.services.CytologyImage.GetCytologyImagesByExternalID(ctx, externalID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &pb.GetCytologyImagesByExternalIdOut{CytologyImages: []*pb.CytologyImage{}}, nil
		}
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	pbImages := make([]*pb.CytologyImage, 0, len(images))
	for _, img := range images {
		pbImages = append(pbImages, mappers.CytologyImageToProto(img))
	}

	return &pb.GetCytologyImagesByExternalIdOut{CytologyImages: pbImages}, nil
}

func (h *handler) GetCytologyImagesByPatientCardId(ctx context.Context, in *pb.GetCytologyImagesByPatientCardIdIn) (*pb.GetCytologyImagesByPatientCardIdOut, error) {
	patientCardID, err := uuid.Parse(in.PatientCardId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "patient_card_id is not a valid uuid: %s", err.Error())
	}

	images, err := h.services.CytologyImage.GetCytologyImagesByPatientCardID(ctx, patientCardID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &pb.GetCytologyImagesByPatientCardIdOut{CytologyImages: []*pb.CytologyImage{}}, nil
		}
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	pbImages := make([]*pb.CytologyImage, 0, len(images))
	for _, img := range images {
		pbImages = append(pbImages, mappers.CytologyImageToProto(img))
	}

	return &pb.GetCytologyImagesByPatientCardIdOut{CytologyImages: pbImages}, nil
}
