package node

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "exam/internal/generated/grpc/service"
	"exam/internal/server/mappers"
	"exam/internal/services/node"
)

func (h *handler) UpdateNode(ctx context.Context, in *pb.UpdateNodeIn) (*pb.UpdateNodeOut, error) {
	if _, err := uuid.Parse(in.Id); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id is not a valid uuid: %s", err.Error())
	}

	node, err := h.services.Node.UpdateNode(
		ctx,
		node.UpdateNodeArg{
			Id:         uuid.MustParse(in.Id),
			Validation: mappers.NodeValidationToDomain(in.Validation),
			Knosp012:   in.Knosp_012,
			Knosp3:     in.Knosp_3,
			Knosp4:     in.Knosp_4,
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	out := new(pb.UpdateNodeOut)
	out.Node = mappers.NodeFromDomain(node)

	return out, nil
}
