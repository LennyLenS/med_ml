package segment

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "exam/internal/generated/grpc/service"
	"exam/internal/services/segment"
)

func (h *handler) CreateSegment(ctx context.Context, in *pb.CreateSegmentIn) (*pb.CreateSegmentOut, error) {
	if _, err := uuid.Parse(in.ImageId); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "image_id is not a valid uuid: %s", err.Error())
	}

	if _, err := uuid.Parse(in.NodeId); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "node_id is not a valid uuid: %s", err.Error())
	}

	if !json.Valid(in.Contor) {
		return nil, status.Errorf(codes.InvalidArgument, "contor is not a valid json")
	}

	id, err := h.services.Segment.CreateManualSegment(ctx, segment.CreateSegmentArg{
		ImageID:  uuid.MustParse(in.ImageId),
		NodeID:   uuid.MustParse(in.NodeId),
		Contor:   in.Contor,
		Knosp012: in.Knosp_012,
		Knosp3:   in.Knosp_3,
		Knosp4:   in.Knosp_4,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	return &pb.CreateSegmentOut{
		Id: id.String(),
	}, nil
}
