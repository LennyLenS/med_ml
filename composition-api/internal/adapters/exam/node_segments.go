package exam

import (
	"context"

	"github.com/google/uuid"

	"composition-api/internal/adapters/exam/mappers"
	domain "composition-api/internal/domain/exam"
	pb "composition-api/internal/generated/grpc/clients/exam"
)

func (a *adapter) CreateNodeWithSegments(ctx context.Context, in CreateNodeWithSegmentsIn) (uuid.UUID, []uuid.UUID, error) {
	req := &pb.CreateNodeWithSegmentsIn{}

	req.MriId = in.MriID.String()

	req.Node = &pb.CreateNodeWithSegmentsIn_Node{
		Knosp_012:   in.Node.Knosp_012,
		Knosp_3:     in.Node.Knosp_3,
		Knosp_4:     in.Node.Knosp_4,
		Description: in.Node.Description,
	}

	for _, segment := range in.Segments {
		req.Segments = append(req.Segments, &pb.CreateNodeWithSegmentsIn_Segment{
			ImageId:   segment.ImageID.String(),
			Contor:    segment.Contor,
			Knosp_012: segment.Knosp_012,
			Knosp_3:   segment.Knosp_3,
			Knosp_4:   segment.Knosp_4,
		})
	}

	res, err := a.client.CreateNodeWithSegments(ctx, req)
	if err != nil {
		return uuid.Nil, nil, err
	}

	segmentIDs := make([]uuid.UUID, 0, len(res.SegmentIds))
	for _, id := range res.SegmentIds {
		segmentIDs = append(segmentIDs, uuid.MustParse(id))
	}

	return uuid.MustParse(res.NodeId), segmentIDs, nil
}

func (a *adapter) GetNodesWithSegmentsByImageId(ctx context.Context, id uuid.UUID) ([]domain.Node, []domain.Segment, error) {
	res, err := a.client.GetNodesWithSegmentsByImageId(ctx, &pb.GetNodesWithSegmentsByImageIdIn{Id: id.String()})
	if err != nil {
		return nil, nil, err
	}

	nodes := mappers.Node{}.SliceDomain(res.Nodes)
	segments := mappers.Segment{}.SliceDomain(res.Segments)

	return nodes, segments, nil
}


