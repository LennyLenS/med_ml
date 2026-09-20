package kt

import (
	"context"

	api "composition-api/internal/generated/http/api"
	services "composition-api/internal/services"
)

type KTHandler interface {
	KtIDGet(ctx context.Context, params api.KtIDGetParams) (api.KtIDGetRes, error)
	KtPost(ctx context.Context, req *api.KtPostReq) (api.KtPostRes, error)
	KtIDPatch(ctx context.Context, request *api.KtIDPatchReq, params api.KtIDPatchParams) (api.KtIDPatchRes, error)
	KtIDDelete(ctx context.Context, params api.KtIDDeleteParams) (api.KtIDDeleteRes, error)
}

type handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) KTHandler {
	return &handler{
		services: services,
	}
}
