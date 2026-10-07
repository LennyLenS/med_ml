package segment

import (
	"context"
	"encoding/json"
	"fmt"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/exam/mappers"
	apimappers "composition-api/internal/server/mappers"
	segmentSrv "composition-api/internal/services/exam/segment"
)

func (h *handler) MriSegmentIDPatch(ctx context.Context, req *api.MriSegmentIDPatchReq, params api.MriSegmentIDPatchParams) (api.MriSegmentIDPatchRes, error) {
	var contor []byte
	if req.Contor != nil {
		contorParsed, err := json.Marshal(req.Contor)
		if err != nil {
			return nil, fmt.Errorf("parse input contor: %w", err)
		}
		contor = contorParsed
	}

	segment, err := h.services.ExamSegmentService.Update(ctx, segmentSrv.UpdateSegmentArg{
		Id:        params.ID,
		Contor:    contor,
		Knosp_012: apimappers.FromOptFloat64(req.Knosp012),
		Knosp_3:   apimappers.FromOptFloat64(req.Knosp3),
		Knosp_4:   apimappers.FromOptFloat64(req.Knosp4),
	})
	if err != nil {
		return nil, err
	}

	resp, err := mappers.Segment{}.Domain(segment)
	if err != nil {
		return nil, fmt.Errorf("map segment to output: %w", err)
	}
	return &resp, nil
}
