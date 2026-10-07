package kt

import (
	"context"

	"github.com/google/uuid"

	pb "exam/internal/generated/grpc/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *handler) DeleteKt(ctx context.Context, req *pb.DeleteKtIn) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid kt id: %v", err)
	}

	err = h.services.Kt.DeleteKt(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete kt: %v", err)
	}

	return &emptypb.Empty{}, nil
}
