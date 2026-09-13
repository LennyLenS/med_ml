package segment

import (
	"context"

	"github.com/AlekSi/pointer"

	api "composition-api/internal/generated/http/api"
	"composition-api/internal/server/exam/mappers"
)

func (h *handler) MriNodesIDSegmentsGet(ctx context.Context, params api.MriNodesIDSegmentsGetParams) (api.MriNodesIDSegmentsGetRes, error) {
	segments, err := h.services.ExamSegmentService.GetByNodeID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	response, err := mappers.Segment{}.SliceDomain(segments)
	if err != nil {
		return nil, err
	}

	return pointer.To(api.MriNodesIDSegmentsGetOKApplicationJSON(response)), nil
}
