package segmentation

import (
	"context"
	"errors"
	"net/http"

	"composition-api/internal/domain"

	"github.com/AlekSi/pointer"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/cytology/mappers"
)

func (h *handler) CytologyIDSegmentationGroupsPost(ctx context.Context, req *api.CytologyIDSegmentationGroupsPostReq, params api.CytologyIDSegmentationGroupsPostParams) (api.CytologyIDSegmentationGroupsPostRes, error) {
	arg := mappers.SegmentationGroup{}.CreateArg(params.ID, req)

	id, err := h.services.CytologyService.CreateSegmentationGroup(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadRequest):
			return &api.CytologyIDSegmentationGroupsPostBadRequest{
				StatusCode: http.StatusBadRequest,
				Response: api.Error{
					Message: "Неверный формат запроса",
				},
			}, nil
		case errors.Is(err, domain.ErrNotFound):
			return &api.CytologyIDSegmentationGroupsPostNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Цитологическое исследование не найдено",
				},
			}, nil
		default:
			return nil, err
		}
	}

	return pointer.To(api.SimpleUuid{ID: id}), nil
}
