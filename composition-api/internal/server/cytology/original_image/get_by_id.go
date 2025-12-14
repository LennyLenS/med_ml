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

func (h *handler) CytologyOriginalImageIDGet(ctx context.Context, params api.CytologyOriginalImageIDGetParams) (api.CytologyOriginalImageIDGetRes, error) {
	img, err := h.services.CytologyService.GetOriginalImageById(ctx, params.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &api.CytologyOriginalImageIDGetNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Оригинальное изображение не найдено",
				},
			}, nil
		}
		return nil, err
	}

	return pointer.To(api.CytologyOriginalImageIDGetOK{
		OriginalImage: api.OptOriginalImage{
			Value: mappers.OriginalImage{}.Domain(img),
			Set:   true,
		},
	}), nil
}
