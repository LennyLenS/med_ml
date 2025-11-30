package segmentation_group

import (
	"context"

	pb "cytology/internal/generated/grpc/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) GetSegmentationGroupsByCytologyId(ctx context.Context, in *pb.GetSegmentationGroupsByCytologyIdIn) (*pb.GetSegmentationGroupsByCytologyIdOut, error) {
	cytologyID, err := uuid.Parse(in.CytologyId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cytology_id is not a valid uuid: %s", err.Error())
	}

	groups, err := h.services.SegmentationGroup.GetSegmentationGroupsByCytologyID(ctx, cytologyID)
	if err != nil {
		return &pb.GetSegmentationGroupsByCytologyIdOut{SegmentationGroups: []*pb.SegmentationGroup{}}, nil
	}

	// TODO: добавить маппинг
	return &pb.GetSegmentationGroupsByCytologyIdOut{SegmentationGroups: []*pb.SegmentationGroup{}}, nil
}
