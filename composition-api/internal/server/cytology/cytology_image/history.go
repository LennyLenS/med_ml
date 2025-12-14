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

func (h *handler) CytologyHistoryRead(ctx context.Context, params api.CytologyHistoryReadParams) (api.CytologyHistoryReadRes, error) {
	id, err := uuid.Parse(params.ID)
	if err != nil {
		return &api.CytologyHistoryReadInternalServerError{
			StatusCode: http.StatusBadRequest,
			Response: api.Error{
				Message: "Неверный формат ID",
			},
		}, nil
	}

	imgs, err := h.services.CytologyService.GetCytologyImageHistory(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &api.CytologyHistoryReadNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Цитологическое исследование не найдено",
				},
			}, nil
		}
		return nil, err
	}

	// Преобразуем в формат согласно swagger.json (пагинированный список CytologyImageModel)
	results := mappers.CytologyImage{}.ToCytologyImageModelList(imgs)

	// Создаем пагинированный ответ
	// TODO: Реализовать правильную пагинацию с limit и offset
	result := api.CytologyHistoryReadOK{
		Count:    len(results),
		Next:     api.OptURI{Set: false},
		Previous: api.OptURI{Set: false},
		Results:  results,
	}

	return pointer.To(result), nil
}
