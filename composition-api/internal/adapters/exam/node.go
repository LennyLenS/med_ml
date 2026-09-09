package exam

import (
	"context"

	"github.com/google/uuid"

	"composition-api/internal/adapters/exam/mappers"
	domain "composition-api/internal/domain/exam"
	pb "composition-api/internal/generated/grpc/clients/exam"
)

var nodeValidationMap = map[domain.NodeValidation]pb.NodeValidation{
	domain.NodeValidationValid:   pb.NodeValidation_NODE_VALIDATION_VALID,
	domain.NodeValidationInvalid: pb.NodeValidation_NODE_VALIDATION_INVALID,
}

func (a *adapter) GetNodesByMriId(ctx context.Context, id uuid.UUID) ([]domain.Node, error) {
	res, err := a.client.GetNodesByMriId(ctx, &pb.GetNodesByMriIdIn{MriId: id.String()})
	if err != nil {
		return nil, err
	}

	return mappers.Node{}.SliceDomain(res.Nodes), nil
}

func (a *adapter) UpdateNode(ctx context.Context, in UpdateNodeIn) (domain.Node, error) {
	res, err := a.client.UpdateNode(ctx, &pb.UpdateNodeIn{
		Id:         in.Id.String(),
		Validation: mappers.PointerFromMap(nodeValidationMap, in.Validation),
		Knosp_012:  in.Knosp_012,
		Knosp_3:    in.Knosp_3,
		Knosp_4:    in.Knosp_4,
	})
	if err != nil {
		return domain.Node{}, err
	}

	return mappers.Node{}.Domain(res.Node), nil
}

func (a *adapter) DeleteNode(ctx context.Context, id uuid.UUID) error {
	_, err := a.client.DeleteNode(ctx, &pb.DeleteNodeIn{Id: id.String()})
	if err != nil {
		return err
	}
	return nil
}


