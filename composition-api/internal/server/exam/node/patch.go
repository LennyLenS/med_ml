package node

import (
	"context"

	"github.com/AlekSi/pointer"

	domain "composition-api/internal/domain/exam"
	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/exam/mappers"
	apimappers "composition-api/internal/server/mappers"
	nodeSrv "composition-api/internal/services/exam/node"
)

func (h *handler) MriNodesIDPatch(ctx context.Context, req *api.MriNodesIDPatchReq, params api.MriNodesIDPatchParams) (api.MriNodesIDPatchRes, error) {
	var validation *domain.NodeValidation
	switch {
	case req.Validation.Null:
		validation = pointer.To(domain.NodeValidationNull)
	case req.Validation.IsSet():
		validation = pointer.To(domain.NodeValidation(req.Validation.Value))
	}

	node, err := h.services.ExamNodeService.UpdateNode(ctx, nodeSrv.UpdateNodeArg{
		Id:         params.ID,
		Validation: validation,
		Knosp_012:  apimappers.FromOptFloat64(req.Knosp012),
		Knosp_3:    apimappers.FromOptFloat64(req.Knosp3),
		Knosp_4:    apimappers.FromOptFloat64(req.Knosp4),
	})
	if err != nil {
		return nil, err
	}

	return pointer.To(mappers.Node{}.Domain(node)), nil
}
