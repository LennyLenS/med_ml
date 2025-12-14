package cytology_image

import (
	"context"
	"errors"
	"net/http"

	"composition-api/internal/domain"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"

	domainCytology "composition-api/internal/domain/cytology"
	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/cytology/mappers"
)

func (h *handler) CytologySegmentsList(ctx context.Context, params api.CytologySegmentsListParams) (api.CytologySegmentsListRes, error) {
	id, err := uuid.Parse(params.ID)
	if err != nil {
		return &api.CytologySegmentsListBadRequest{
			StatusCode: http.StatusBadRequest,
			Response: api.Error{
				Message: "Неверный формат ID",
			},
		}, nil
	}

	var segType *domainCytology.SegType
	if params.SegType.Set {
		st := domainCytology.SegType(params.SegType.Value)
		segType = &st
	}

	var groupType *domainCytology.GroupType
	if params.GroupType.Set {
		gt := domainCytology.GroupType(params.GroupType.Value)
		groupType = &gt
	}

	var isAI *bool
	if params.IsAI.Set {
		isAI = &params.IsAI.Value
	}

	groups, err := h.services.CytologyService.GetSegmentationGroupsByCytologyId(ctx, id, segType, groupType, isAI)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &api.CytologySegmentsListNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Цитологическое исследование не найдено",
				},
			}, nil
		}
		return nil, err
	}

	// Преобразуем в формат согласно swagger.json (пагинированный список SegmentationData)
	results := mappers.SegmentationGroup{}.ToSegmentationDataList(groups)

	// Создаем пагинированный ответ
	// TODO: Реализовать правильную пагинацию с limit и offset
	result := api.CytologySegmentsListOK{
		Count:    len(results),
		Next:     api.OptURI{Set: false},
		Previous: api.OptURI{Set: false},
		Results:  results,
	}

	return pointer.To(result), nil
}
