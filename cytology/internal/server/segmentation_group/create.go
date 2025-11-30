package segmentation_group

import (
	"context"

	pb "cytology/internal/generated/grpc/service"
	"cytology/internal/domain"
	"cytology/internal/services/segmentation_group"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) CreateSegmentationGroup(ctx context.Context, in *pb.CreateSegmentationGroupIn) (*pb.CreateSegmentationGroupOut, error) {
	cytologyID, err := uuid.Parse(in.CytologyId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cytology_id is not a valid uuid: %s", err.Error())
	}

	var details []byte
	if in.Details != nil && *in.Details != "" {
		details = []byte(*in.Details)
	}

	arg := segmentation_group.CreateSegmentationGroupArg{
		CytologyID: cytologyID,
		SegType:    domain.SegType(in.SegType.String()),
		GroupType:  domain.GroupType(in.GroupType.String()),
		IsAI:       in.IsAi,
		Details:    details,
	}

	id, err := h.services.SegmentationGroup.CreateSegmentationGroup(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	return &pb.CreateSegmentationGroupOut{Id: id.String()}, nil
}
