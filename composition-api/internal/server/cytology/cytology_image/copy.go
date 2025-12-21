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
	// Пока используем pk или id как строку UUID, если они переданы
	var id uuid.UUID
	var err error

	if req.ID.Set {
		// Пытаемся распарсить как UUID (если передана строка UUID)
		id, err = uuid.Parse(string(rune(req.ID.Value)))
		if err != nil {
			// Если не UUID, то это integer ID - нужно получить UUID по integer ID
			// TODO: Реализовать получение UUID по integer ID через другой сервис или БД
			return &api.CytologyCopyCreateBadRequest{
				StatusCode: http.StatusBadRequest,
				Response: api.Error{
					Message: "Неверный формат ID. Ожидается UUID",
				},
			}, nil
		}
	} else if req.Pk.Set {
		// Пытаемся распарсить как UUID
		id, err = uuid.Parse(string(rune(req.Pk.Value)))
		if err != nil {
			return &api.CytologyCopyCreateBadRequest{
				StatusCode: http.StatusBadRequest,
				Response: api.Error{
					Message: "Неверный формат pk. Ожидается UUID",
				},
			}, nil
		}
	} else {
		return &api.CytologyCopyCreateBadRequest{
			StatusCode: http.StatusBadRequest,
			Response: api.Error{
				Message: "ID или pk обязательны",
			},
		}, nil
	}

	// Копируем исследование
	_, err = h.services.CytologyService.CopyCytologyImage(ctx, id)
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
	// UUID не преобразуется в integer, оставляем пустым
	result := api.CytologyCopyCreateCreated{
		Pk: api.OptInt{
			Set: false,
		},
		ID: api.OptInt{
			Set: false,
		},
	}

	return pointer.To(result), nil
}
