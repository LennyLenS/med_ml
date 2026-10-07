package node_segment

import (
	"context"
	"encoding/json"

	api "composition-api/internal/generated/http/api"
	"composition-api/internal/server/mappers"
	"composition-api/internal/services/exam/node_segment"
)

func (h *handler) MriIDNodesSegmentsPost(ctx context.Context, req *api.MriIDNodesSegmentsPostReq, params api.MriIDNodesSegmentsPostParams) (api.MriIDNodesSegmentsPostRes, error) {
	arg := node_segment.CreateNodeWithSegmentArg{}
	arg.MriID = params.ID
	arg.Node = node_segment.CreateNodeWithSegmentArg_Node{
		Knosp_012:   req.Node.Knosp012,
		Knosp_3:     req.Node.Knosp3,
		Knosp_4:     req.Node.Knosp4,
		Description: mappers.FromOptString(req.Node.Description),
	}

	for _, segment := range req.Segments {
		contor, err := json.Marshal(segment.Contor)
		if err != nil {
			return nil, err
		}

		arg.Segments = append(arg.Segments, node_segment.CreateNodeWithSegmentArg_Segment{
			ImageID:   segment.ImageID,
			Contor:    contor,
			Knosp_012: segment.Knosp012,
			Knosp_3:   segment.Knosp3,
			Knosp_4:   segment.Knosp4,
		})
	}

	nodeID, segmentIDs, err := h.services.ExamNodeSegmentService.CreateNodeWithSegment(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &api.MriIDNodesSegmentsPostOK{
		NodeID:     nodeID,
		SegmentIds: segmentIDs,
	}, nil
}
