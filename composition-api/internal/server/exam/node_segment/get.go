package node_segment

import (
	"context"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/exam/mappers"
)

func (h *handler) MriImageIDNodesSegmentsGet(ctx context.Context, params api.MriImageIDNodesSegmentsGetParams) (api.MriImageIDNodesSegmentsGetRes, error) {
	nodes, segments, err := h.services.ExamNodeSegmentService.GetNodeWithSegmentsByImageID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	segmentsResp, err := mappers.Segment{}.SliceDomain(segments)
	if err != nil {
		return nil, err
	}

	return &api.MriImageIDNodesSegmentsGetOK{
		Nodes:    mappers.Node{}.SliceDomain(nodes),
		Segments: segmentsResp,
	}, nil
}
