package cytology_image

import (
	"context"
	"errors"
	"net/http"

	"composition-api/internal/domain"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"

	api "composition-api/internal/generated/http/api"
)

func (h *handler) CytologyCopyCreate(ctx context.Context, req *api.CytologyCopyCreateReq) (api.CytologyCopyCreateRes, error) {
	// В swagger.json req содержит pk и id (оба integer)
	// Нужно преобразовать в UUID для вызова сервиса
	var id uuid.UUID
	if req.ID.Set {
		// TODO: Преобразовать integer ID в UUID
		// Пока используем pk, если он есть
		if req.Pk.Set {
			// Временная заглушка - нужно будет реализовать правильное преобразование
			id = uuid.New()
		}
	} else if req.Pk.Set {
		// Временная заглушка
		id = uuid.New()
	} else {
		return &api.CytologyCopyCreateBadRequest{
			StatusCode: http.StatusBadRequest,
			Response: api.Error{
				Message: "ID или pk обязательны",
			},
		}, nil
	}

	img, err := h.services.CytologyService.CopyCytologyImage(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &api.CytologyCopyCreateNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Цитологическое исследование не найдено",
				},
			}, nil
		case errors.Is(err, domain.ErrBadRequest):
			return &api.CytologyCopyCreateBadRequest{
				StatusCode: http.StatusBadRequest,
				Response: api.Error{
					Message: "Неверный формат запроса",
				},
			}, nil
		default:
			return nil, err
		}
	}

	// Возвращаем согласно swagger.json (CytologyImageCopy с pk и id)
	result := api.CytologyCopyCreateCreated{
		Pk: api.OptInt{
			// TODO: Преобразовать UUID в integer
			Set: false,
		},
		ID: api.OptInt{
			// TODO: Преобразовать UUID в integer
			Set: false,
		},
	}

	return pointer.To(result), nil
}
