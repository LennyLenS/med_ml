package kt

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "exam/internal/generated/grpc/service"
	"exam/internal/services/kt"
)

func (h *handler) CreateKt(ctx context.Context, in *pb.CreateKtIn) (*pb.CreateKtOut, error) {
	if _, err := uuid.Parse(in.Author); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "author is not a valid uuid: %s", err.Error())
	}

	id, err := h.services.Kt.CreateKt(ctx, kt.CreateKTArg{
		Author:      uuid.MustParse(in.Author),
		DeviceID:    int(in.DeviceId),
		Description: in.Description,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	return &pb.CreateKtOut{Id: id.String()}, nil
}
