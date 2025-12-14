package cytology_image

import (
	"context"
	"errors"
	"net/http"

	"composition-api/internal/domain"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/cytology/mappers"
)

func (h *handler) CytologyRead(ctx context.Context, params api.CytologyReadParams) (api.CytologyReadRes, error) {
	id, err := uuid.Parse(params.ID)
	if err != nil {
		return &api.CytologyReadInternalServerError{
			StatusCode: http.StatusBadRequest,
			Response: api.Error{
				Message: "Неверный формат ID",
			},
		}, nil
	}

	img, origImg, err := h.services.CytologyService.GetCytologyImageById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &api.CytologyReadNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Цитологическое исследование не найдено",
				},
			}, nil
		default:
			return nil, err
		}
	}

	// Маппим в структуру согласно swagger.json
	result := api.CytologyReadOK{
		OriginalImage: mappers.OriginalImage{}.ToCytologyReadOKOriginalImage(origImg),
		Info:          mappers.CytologyImage{}.ToCytologyReadOKInfo(img),
	}

	return pointer.To(result), nil
}
