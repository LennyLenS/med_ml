package mri

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "exam/internal/generated/grpc/service"
	"exam/internal/server/mappers"
	"exam/internal/services/mri"
)

func (h *handler) CreateMri(ctx context.Context, in *pb.CreateMriIn) (*pb.CreateMriOut, error) {
	if _, err := uuid.Parse(in.ExternalId); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "external_id is not a valid uuid: %s", err.Error())
	}

	if _, err := uuid.Parse(in.Author); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "author is not a valid uuid: %s", err.Error())
	}

	id, err := h.services.Mri.CreateMri(ctx, mri.CreateMriArg{
		Projection:  mappers.MriProjectionReverseMap[in.Projection],
		ExternalID:  uuid.MustParse(in.ExternalId),
		Author:      uuid.MustParse(in.Author),
		DeviceID:    int(in.DeviceId),
		Description: in.Description,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	return &pb.CreateMriOut{Id: id.String()}, nil
}
