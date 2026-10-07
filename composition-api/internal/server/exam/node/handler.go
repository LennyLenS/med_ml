package node

import (
	"context"

	api "composition-api/internal/generated/http/api"
	services "composition-api/internal/services"
)

type NodeHandler interface {
	MriIDNodesGet(ctx context.Context, params api.MriIDNodesGetParams) (api.MriIDNodesGetRes, error)
	MriNodesIDPatch(ctx context.Context, req *api.MriNodesIDPatchReq, params api.MriNodesIDPatchParams) (api.MriNodesIDPatchRes, error)
	MriNodesIDDelete(ctx context.Context, params api.MriNodesIDDeleteParams) (api.MriNodesIDDeleteRes, error)
}

type handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) NodeHandler {
	return &handler{
		services: services,
	}
}
