package cytology_image

import (
	"context"
	"errors"
	"net/http"

	"composition-api/internal/domain"

	"github.com/AlekSi/pointer"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/cytology/mappers"
)

func (h *handler) CytologyPost(ctx context.Context, req *api.CytologyPostReq) (api.CytologyPostRes, error) {
	arg := mappers.CytologyImage{}.CreateArg(req)

	id, err := h.services.CytologyService.CreateCytologyImage(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadRequest):
			return &api.CytologyPostBadRequest{
				StatusCode: http.StatusBadRequest,
				Response: api.Error{
					Message: "Неверный формат запроса",
				},
			}, nil
		case errors.Is(err, domain.ErrUnprocessableEntity):
			return &api.CytologyPostUnprocessableEntity{
				StatusCode: http.StatusUnprocessableEntity,
				Response: api.Error{
					Message: "Ошибка валидации данных",
				},
			}, nil
		default:
			return nil, err
		}
	}

	return pointer.To(api.SimpleUuid{ID: id}), nil
}
