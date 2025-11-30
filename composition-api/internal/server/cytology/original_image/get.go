package original_image

import (
	"context"
	"errors"
	"net/http"

	"composition-api/internal/domain"

	"github.com/AlekSi/pointer"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/cytology/mappers"
)

func (h *handler) CytologyIDOriginalImageGet(ctx context.Context, params api.CytologyIDOriginalImageGetParams) (api.CytologyIDOriginalImageGetRes, error) {
	imgs, err := h.services.CytologyService.GetOriginalImagesByCytologyId(ctx, params.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &api.CytologyIDOriginalImageGetNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Оригинальные изображения не найдены",
				},
			}, nil
		}
		return nil, err
	}

	return pointer.To(api.CytologyIDOriginalImageGetOKApplicationJSON(
		mappers.OriginalImage{}.SliceDomain(imgs),
	)), nil
}
