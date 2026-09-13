package mri

import (
	"context"

	pb "exam/internal/generated/grpc/service"
	"exam/internal/services"

	"google.golang.org/protobuf/types/known/emptypb"
)

type MriHandler interface {
	CreateMri(ctx context.Context, req *pb.CreateMriIn) (*pb.CreateMriOut, error)

	GetMriById(ctx context.Context, req *pb.GetMriByIdIn) (*pb.GetMriByIdOut, error)
	GetMrisByExternalId(ctx context.Context, req *pb.GetMrisByExternalIdIn) (*pb.GetMrisByExternalIdOut, error)
	GetMrisByAuthor(ctx context.Context, req *pb.GetMrisByAuthorIn) (*pb.GetMrisByAuthorOut, error)
	GetEchographicByMriId(ctx context.Context, req *pb.GetEchographicByMriIdIn) (*pb.GetEchographicByMriIdOut, error)

	UpdateMri(ctx context.Context, req *pb.UpdateMriIn) (*pb.UpdateMriOut, error)
	UpdateEchographic(ctx context.Context, in *pb.UpdateEchographicIn) (*pb.UpdateEchographicOut, error)

	DeleteMri(ctx context.Context, req *pb.DeleteMriIn) (*emptypb.Empty, error)
}

type handler struct {
	services *services.Services
}

func New(
	services *services.Services,
) MriHandler {
	return &handler{
		services: services,
	}
}
