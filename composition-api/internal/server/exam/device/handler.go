package device

import (
	"context"

	api "composition-api/internal/generated/http/api"
	services "composition-api/internal/services"
)

type DeviceHandler interface {
	MriDevicePost(ctx context.Context, req *api.MriDevicePostReq) (api.MriDevicePostRes, error)
	MriDevicesGet(ctx context.Context) (api.MriDevicesGetRes, error)
}

type handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) DeviceHandler {
	return &handler{
		services: services,
	}
}
