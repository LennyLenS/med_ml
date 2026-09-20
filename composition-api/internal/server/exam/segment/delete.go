package segment

import (
	"context"

	api "composition-api/internal/generated/http/api"
)

func (h *handler) MriSegmentIDDelete(ctx context.Context, params api.MriSegmentIDDeleteParams) (api.MriSegmentIDDeleteRes, error) {
	err := h.services.ExamSegmentService.Delete(ctx, params.ID)
	if err != nil {
		return nil, err
	}
	return &api.MriSegmentIDDeleteOK{}, nil
}
