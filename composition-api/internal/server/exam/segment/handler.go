package segment

import (
	"context"

	api "composition-api/internal/generated/http/api"
	services "composition-api/internal/services"
)

type SegmentHandler interface {
	MriSegmentPost(ctx context.Context, req *api.MriSegmentPostReq) (api.MriSegmentPostRes, error)
	MriNodesIDSegmentsGet(ctx context.Context, params api.MriNodesIDSegmentsGetParams) (api.MriNodesIDSegmentsGetRes, error)
	MriSegmentIDPatch(ctx context.Context, req *api.MriSegmentIDPatchReq, params api.MriSegmentIDPatchParams) (api.MriSegmentIDPatchRes, error)
	MriSegmentIDDelete(ctx context.Context, params api.MriSegmentIDDeleteParams) (api.MriSegmentIDDeleteRes, error)
}

type handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) SegmentHandler {
	return &handler{
		services: services,
	}
}
