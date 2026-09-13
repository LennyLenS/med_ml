package node

import (
	"context"

	api "composition-api/internal/generated/http/api"
)

func (h *handler) MriNodesIDDelete(ctx context.Context, params api.MriNodesIDDeleteParams) (api.MriNodesIDDeleteRes, error) {
	err := h.services.ExamNodeService.DeleteNode(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.MriNodesIDDeleteOK{}, nil
}
