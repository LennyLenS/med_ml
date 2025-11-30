package original_image

import (
	"context"

	pb "cytology/internal/generated/grpc/service"
	"cytology/internal/services/original_image"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) CreateOriginalImage(ctx context.Context, in *pb.CreateOriginalImageIn) (*pb.CreateOriginalImageOut, error) {
	cytologyID, err := uuid.Parse(in.CytologyId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cytology_id is not a valid uuid: %s", err.Error())
	}

	arg := original_image.CreateOriginalImageArg{
		CytologyID: cytologyID,
		ImagePath:  in.ImagePath,
		DelayTime:  in.DelayTime,
	}

	id, err := h.services.OriginalImage.CreateOriginalImage(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	return &pb.CreateOriginalImageOut{Id: id.String()}, nil
}
