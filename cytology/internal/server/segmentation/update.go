package segmentation

import (
	"context"

	pb "cytology/internal/generated/grpc/service"
	"cytology/internal/domain"
	"cytology/internal/server/mappers"
	"cytology/internal/services/segmentation"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) UpdateSegmentation(ctx context.Context, in *pb.UpdateSegmentationIn) (*pb.UpdateSegmentationOut, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id is not a valid uuid: %s", err.Error())
	}

	points := make([]domain.SegmentationPoint, 0, len(in.Points))
	for _, p := range in.Points {
		points = append(points, domain.SegmentationPoint{
			X: int(p.X),
			Y: int(p.Y),
		})
	}

	arg := segmentation.UpdateSegmentationArg{
		Id:     id,
		Points: points,
	}

	seg, err := h.services.Segmentation.UpdateSegmentation(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	return &pb.UpdateSegmentationOut{
		Segmentation: mappers.SegmentationToProto(seg),
	}, nil
}
