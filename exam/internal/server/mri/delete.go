package mri

import (
	"context"

	"github.com/google/uuid"

	pb "exam/internal/generated/grpc/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *handler) DeleteMri(ctx context.Context, req *pb.DeleteMriIn) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid mri id: %v", err)
	}

	err = h.services.Mri.DeleteMri(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete mri: %v", err)
	}

	return &emptypb.Empty{}, nil
}
