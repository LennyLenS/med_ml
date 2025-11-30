package cytology_image

import (
	"context"
	"errors"
	"net/http"

	"composition-api/internal/domain"

	api "composition-api/internal/generated/http/api"
)

func (h *handler) CytologyIDDelete(ctx context.Context, params api.CytologyIDDeleteParams) (api.CytologyIDDeleteRes, error) {
	err := h.services.CytologyService.DeleteCytologyImage(ctx, params.ID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &api.CytologyIDDeleteNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Цитологическое исследование не найдено",
				},
			}, nil
		default:
			return nil, err
		}
	}

	return &api.CytologyIDDeleteOK{}, nil
}
