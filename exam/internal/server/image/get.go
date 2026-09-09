package image

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "exam/internal/generated/grpc/service"
	"exam/internal/server/mappers"
)

func (h *handler) GetImagesByMriId(ctx context.Context, in *pb.GetImagesByMriIdIn) (*pb.GetImagesByMriIdOut, error) {
	if _, err := uuid.Parse(in.MriId); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "mri_id is not a valid uuid: %s", err.Error())
	}

	images, err := h.services.Image.GetImagesByMriID(ctx, uuid.MustParse(in.MriId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	out := new(pb.GetImagesByMriIdOut)
	out.Images = mappers.SliceImageFromDomain(images)

	return out, nil
}
