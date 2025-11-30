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

func (h *handler) CytologySegmentationGroupIDSegmentsPost(ctx context.Context, req *api.CytologySegmentationGroupIDSegmentsPostReq, params api.CytologySegmentationGroupIDSegmentsPostParams) (api.CytologySegmentationGroupIDSegmentsPostRes, error) {
	arg := mappers.Segmentation{}.CreateArg(params.ID, req)

	id, err := h.services.CytologyService.CreateSegmentation(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadRequest):
			return &api.CytologySegmentationGroupIDSegmentsPostBadRequest{
				StatusCode: http.StatusBadRequest,
				Response: api.Error{
					Message: "Неверный формат запроса",
				},
			}, nil
		case errors.Is(err, domain.ErrNotFound):
			return &api.CytologySegmentationGroupIDSegmentsPostNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Группа сегментаций не найдена",
				},
			}, nil
		default:
			return nil, err
		}
	}

	return pointer.To(api.SimpleUuid{ID: id}), nil
}
