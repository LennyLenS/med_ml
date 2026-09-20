package image

import (
	"context"

	api "composition-api/internal/generated/http/api"
	services "composition-api/internal/services"
)

type ImageHandler interface {
	MriIDImagesGet(ctx context.Context, params api.MriIDImagesGetParams) (api.MriIDImagesGetRes, error)
}

type handler struct {
	services *services.Services
}

func NewHandler(services *services.Services) ImageHandler {
	return &handler{
		services: services,
	}
}
