package node

import (
	"context"

	"github.com/AlekSi/pointer"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/exam/mappers"
)

func (h *handler) MriIDNodesGet(ctx context.Context, params api.MriIDNodesGetParams) (api.MriIDNodesGetRes, error) {
	nodes, err := h.services.ExamNodeService.GetNodesByMriID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return pointer.To(api.MriIDNodesGetOKApplicationJSON(mappers.Node{}.SliceDomain(nodes))), nil
}
