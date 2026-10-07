package image

import (
	"context"

	pb "exam/internal/generated/grpc/service"
	"exam/internal/services"
)

type ImageHandler interface {
	GetImagesByMriId(ctx context.Context, in *pb.GetImagesByMriIdIn) (*pb.GetImagesByMriIdOut, error)
}

type handler struct {
	services *services.Services
}

func New(
	services *services.Services,
) ImageHandler {
	return &handler{
		services: services,
	}
}
