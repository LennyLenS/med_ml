package node_segment

import (
	"context"

	api "composition-api/internal/generated/http/api"
	services "composition-api/internal/services"
)

type NodeSegmentHandler interface {
	MriIDNodesSegmentsPost(ctx context.Context, req *api.MriIDNodesSegmentsPostReq, params api.MriIDNodesSegmentsPostParams) (api.MriIDNodesSegmentsPostRes, error)
	MriImageIDNodesSegmentsGet(ctx context.Context, params api.MriImageIDNodesSegmentsGetParams) (api.MriImageIDNodesSegmentsGetRes, error)
}

type handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) NodeSegmentHandler {
	return &handler{
		services: services,
	}
}
