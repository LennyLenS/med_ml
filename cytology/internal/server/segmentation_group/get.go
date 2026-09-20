package segmentation_group

import (
	"context"

	"cytology/internal/domain"
	pb "cytology/internal/generated/grpc/service"
	"cytology/internal/server/mappers"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) GetSegmentationGroupsByCytologyId(ctx context.Context, in *pb.GetSegmentationGroupsByCytologyIdIn) (*pb.GetSegmentationGroupsByCytologyIdOut, error) {
	cytologyID, err := uuid.Parse(in.CytologyId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cytology_id is not a valid uuid: %s", err.Error())
	}

	var segType *domain.SegType
	if in.SegType != nil {
		value := mappers.SegTypeReverseMap[in.GetSegType()]
		segType = &value
	}
	var groupType *domain.GroupType
	if in.GroupType != nil {
		value := mappers.GroupTypeReverseMap[in.GetGroupType()]
		groupType = &value
	}

	groups, err := h.services.SegmentationGroup.GetSegmentationGroupsByCytologyID(ctx, cytologyID, segType, groupType, in.IsAi)
	if err != nil {
		return &pb.GetSegmentationGroupsByCytologyIdOut{SegmentationGroups: []*pb.SegmentationGroup{}}, nil
	}

	return &pb.GetSegmentationGroupsByCytologyIdOut{
		SegmentationGroups: mappers.SegmentationGroupSliceToProto(groups),
	}, nil
}
