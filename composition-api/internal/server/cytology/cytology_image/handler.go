package cytology_image

import (
	"context"

	api "composition-api/internal/generated/http/api"
	services "composition-api/internal/services"
)

type CytologyImageHandler interface {
	CytologyPost(ctx context.Context, req *api.CytologyPostReq) (api.CytologyPostRes, error)
	CytologyIDGet(ctx context.Context, params api.CytologyIDGetParams) (api.CytologyIDGetRes, error)
	CytologyIDPatch(ctx context.Context, req *api.CytologyIDPatchReq, params api.CytologyIDPatchParams) (api.CytologyIDPatchRes, error)
	CytologyIDDelete(ctx context.Context, params api.CytologyIDDeleteParams) (api.CytologyIDDeleteRes, error)
	CytologiesExternalIDGet(ctx context.Context, params api.CytologiesExternalIDGetParams) (api.CytologiesExternalIDGetRes, error)
	CytologiesPatientCardDoctorIDPatientIDGet(ctx context.Context, params api.CytologiesPatientCardDoctorIDPatientIDGetParams) (api.CytologiesPatientCardDoctorIDPatientIDGetRes, error)
}

type handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) CytologyImageHandler {
	return &handler{
		services: services,
	}
}
