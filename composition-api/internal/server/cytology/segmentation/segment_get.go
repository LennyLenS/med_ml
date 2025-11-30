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

func (h *handler) CytologySegmentationGroupIDSegmentsGet(ctx context.Context, params api.CytologySegmentationGroupIDSegmentsGetParams) (api.CytologySegmentationGroupIDSegmentsGetRes, error) {
	segments, err := h.services.CytologyService.GetSegmentsByGroupId(ctx, params.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &api.CytologySegmentationGroupIDSegmentsGetNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Сегменты не найдены",
				},
			}, nil
		}
		return nil, err
	}

	return pointer.To(api.CytologySegmentationGroupIDSegmentsGetOKApplicationJSON(
		mappers.Segmentation{}.SliceDomain(segments),
	)), nil
}
