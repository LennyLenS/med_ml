package node

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "exam/internal/generated/grpc/service"
	"exam/internal/services"
)

type NodeHandler interface {
	GetNodesByMriId(ctx context.Context, in *pb.GetNodesByMriIdIn) (*pb.GetNodesByMriIdOut, error)
	UpdateNode(ctx context.Context, in *pb.UpdateNodeIn) (*pb.UpdateNodeOut, error)
	DeleteNode(ctx context.Context, in *pb.DeleteNodeIn) (*empty.Empty, error)
}

type handler struct {
	services *services.Services
}

func New(
	services *services.Services,
) NodeHandler {
	return &handler{
		services: services,
	}
}
