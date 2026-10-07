package kt

import (
	"context"

	pb "exam/internal/generated/grpc/service"
	"exam/internal/services"

	"google.golang.org/protobuf/types/known/emptypb"
)

type KtHandler interface {
	CreateKt(ctx context.Context, req *pb.CreateKtIn) (*pb.CreateKtOut, error)

	GetKtById(ctx context.Context, req *pb.GetKtByIdIn) (*pb.GetKtByIdOut, error)

	GetKtsByAuthor(ctx context.Context, req *pb.GetKtsByAuthorIn) (*pb.GetKtsByAuthorOut, error)

	UpdateKt(ctx context.Context, in *pb.UpdateKtIn) (*pb.UpdateKtOut, error)

	DeleteKt(ctx context.Context, req *pb.DeleteKtIn) (*emptypb.Empty, error)
}

type handler struct {
	services *services.Services
}

func New(
	services *services.Services,
) KtHandler {
	return &handler{
		services: services,
	}
}
