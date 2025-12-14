package original_image

import (
	"context"

	api "composition-api/internal/generated/http/api"
	services "composition-api/internal/services"
)

type OriginalImageHandler interface {
	CytologyIDOriginalImagePost(ctx context.Context, req *api.CytologyIDOriginalImagePostReq, params api.CytologyIDOriginalImagePostParams) (api.CytologyIDOriginalImagePostRes, error)
	CytologyIDOriginalImageGet(ctx context.Context, params api.CytologyIDOriginalImageGetParams) (api.CytologyIDOriginalImageGetRes, error)
	CytologyOriginalImageIDGet(ctx context.Context, params api.CytologyOriginalImageIDGetParams) (api.CytologyOriginalImageIDGetRes, error)
	CytologyOriginalImageIDPatch(ctx context.Context, req *api.CytologyOriginalImageIDPatchReq, params api.CytologyOriginalImageIDPatchParams) (api.CytologyOriginalImageIDPatchRes, error)
}

type handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) OriginalImageHandler {
	return &handler{
		services: services,
	}
}
