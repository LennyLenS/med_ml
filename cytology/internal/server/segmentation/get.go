package segmentation

import (
	"context"

	pb "cytology/internal/generated/grpc/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) GetSegmentationById(ctx context.Context, in *pb.GetSegmentationByIdIn) (*pb.GetSegmentationByIdOut, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id is not a valid uuid: %s", err.Error())
	}

	seg, err := h.services.Segmentation.GetSegmentationByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	// TODO: добавить маппинг
	_ = seg
	return &pb.GetSegmentationByIdOut{}, nil
}

func (h *handler) GetSegmentsByGroupId(ctx context.Context, in *pb.GetSegmentsByGroupIdIn) (*pb.GetSegmentsByGroupIdOut, error) {
	groupID, err := uuid.Parse(in.SegmentationGroupId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "segmentation_group_id is not a valid uuid: %s", err.Error())
	}

	segs, err := h.services.Segmentation.GetSegmentsByGroupID(ctx, groupID)
	if err != nil {
		return &pb.GetSegmentsByGroupIdOut{Segmentations: []*pb.Segmentation{}}, nil
	}

	// TODO: добавить маппинг
	_ = segs
	return &pb.GetSegmentsByGroupIdOut{Segmentations: []*pb.Segmentation{}}, nil
}
