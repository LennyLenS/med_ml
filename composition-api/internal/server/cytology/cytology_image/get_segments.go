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

func (h *handler) CytologyIDSegmentsGet(ctx context.Context, params api.CytologyIDSegmentsGetParams) (api.CytologyIDSegmentsGetRes, error) {
	groups, err := h.services.CytologyService.GetSegmentationGroupsByCytologyId(ctx, params.ID, nil, nil, nil)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &api.CytologyIDSegmentsGetNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Цитологическое исследование не найдено",
				},
			}, nil
		}
		return nil, err
	}

	// Получаем группы сегментаций
	result := make([]api.SegmentationGroup, 0, len(groups))
	for _, group := range groups {
		groupProto := mappers.SegmentationGroup{}.Domain(group)
		result = append(result, groupProto)
	}

	return pointer.To(api.CytologyIDSegmentsGetOKApplicationJSON(result)), nil
}
