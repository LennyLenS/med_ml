package segment

import (
	"context"
	"encoding/json"

	api "composition-api/internal/generated/http/api"
	segmentSrv "composition-api/internal/services/exam/segment"
)

func (h *handler) MriSegmentPost(ctx context.Context, req *api.MriSegmentPostReq) (api.MriSegmentPostRes, error) {
	contor, err := json.Marshal(req.Contor)
	if err != nil {
		return nil, err
	}

	segmentID, err := h.services.ExamSegmentService.Create(ctx, segmentSrv.CreateSegmentArg{
		ImageID:   req.ImageID,
		NodeID:    req.NodeID,
		Contor:    contor,
		Knosp_012: req.Knosp012,
		Knosp_3:   req.Knosp3,
		Knosp_4:   req.Knosp4,
	})
	if err != nil {
		return nil, err
	}
	return &api.SimpleUuid{ID: segmentID}, nil
}
