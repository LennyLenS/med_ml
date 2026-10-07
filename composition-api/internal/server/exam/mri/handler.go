package mri

import (
	"context"

	api "composition-api/internal/generated/http/api"
	services "composition-api/internal/services"
)

type MriHandler interface {
	MriPost(ctx context.Context, req *api.MriPostReq) (api.MriPostRes, error)
	MriIDGet(ctx context.Context, params api.MriIDGetParams) (api.MriIDGetRes, error)
	MrisExternalIDGet(ctx context.Context, params api.MrisExternalIDGetParams) (api.MrisExternalIDGetRes, error)
	MrisAuthorIDGet(ctx context.Context, params api.MrisAuthorIDGetParams) (api.MrisAuthorIDGetRes, error)
	MriIDEchographicsGet(ctx context.Context, params api.MriIDEchographicsGetParams) (api.MriIDEchographicsGetRes, error)
	MriIDPatch(ctx context.Context, req *api.MriIDPatchReq, params api.MriIDPatchParams) (api.MriIDPatchRes, error)
	MriIDEchographicsPatch(ctx context.Context, req *api.ExamEchographics, params api.MriIDEchographicsPatchParams) (api.MriIDEchographicsPatchRes, error)
	MriIDDelete(ctx context.Context, params api.MriIDDeleteParams) (api.MriIDDeleteRes, error)
}

type handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) MriHandler {
	return &handler{
		services: services,
	}
}
