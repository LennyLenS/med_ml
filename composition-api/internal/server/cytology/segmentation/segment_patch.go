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

func (h *handler) CytologySegmentationIDPatch(ctx context.Context, req *api.CytologySegmentationIDPatchReq, params api.CytologySegmentationIDPatchParams) (api.CytologySegmentationIDPatchRes, error) {
	arg := mappers.Segmentation{}.UpdateArg(params.ID, req)

	segment, err := h.services.CytologyService.UpdateSegmentation(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &api.CytologySegmentationIDPatchNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Сегмент не найден",
				},
			}, nil
		default:
			return nil, err
		}
	}

	return pointer.To(mappers.Segmentation{}.Domain(segment)), nil
}
