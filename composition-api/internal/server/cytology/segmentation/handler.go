package segmentation

import (
	"context"

	api "composition-api/internal/generated/http/api"
	services "composition-api/internal/services"
)

type SegmentationHandler interface {
	CytologyIDSegmentationGroupsPost(ctx context.Context, req *api.CytologyIDSegmentationGroupsPostReq, params api.CytologyIDSegmentationGroupsPostParams) (api.CytologyIDSegmentationGroupsPostRes, error)
	CytologyIDSegmentationGroupsGet(ctx context.Context, params api.CytologyIDSegmentationGroupsGetParams) (api.CytologyIDSegmentationGroupsGetRes, error)
	CytologySegmentationGroupIDPatch(ctx context.Context, req *api.CytologySegmentationGroupIDPatchReq, params api.CytologySegmentationGroupIDPatchParams) (api.CytologySegmentationGroupIDPatchRes, error)
	CytologySegmentationGroupIDDelete(ctx context.Context, params api.CytologySegmentationGroupIDDeleteParams) (api.CytologySegmentationGroupIDDeleteRes, error)
	CytologySegmentationGroupIDSegmentsPost(ctx context.Context, req *api.CytologySegmentationGroupIDSegmentsPostReq, params api.CytologySegmentationGroupIDSegmentsPostParams) (api.CytologySegmentationGroupIDSegmentsPostRes, error)
	CytologySegmentationGroupIDSegmentsGet(ctx context.Context, params api.CytologySegmentationGroupIDSegmentsGetParams) (api.CytologySegmentationGroupIDSegmentsGetRes, error)
	CytologySegmentationIDPatch(ctx context.Context, req *api.CytologySegmentationIDPatchReq, params api.CytologySegmentationIDPatchParams) (api.CytologySegmentationIDPatchRes, error)
	CytologySegmentationIDDelete(ctx context.Context, params api.CytologySegmentationIDDeleteParams) (api.CytologySegmentationIDDeleteRes, error)
}

type handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) SegmentationHandler {
	return &handler{
		services: services,
	}
}
