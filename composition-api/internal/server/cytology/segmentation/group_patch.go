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

func (h *handler) CytologySegmentationGroupIDPatch(ctx context.Context, req *api.CytologySegmentationGroupIDPatchReq, params api.CytologySegmentationGroupIDPatchParams) (api.CytologySegmentationGroupIDPatchRes, error) {
	arg := mappers.SegmentationGroup{}.UpdateArg(params.ID, req)

	group, err := h.services.CytologyService.UpdateSegmentationGroup(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &api.CytologySegmentationGroupIDPatchNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Группа сегментаций не найдена",
				},
			}, nil
		default:
			return nil, err
		}
	}

	return pointer.To(mappers.SegmentationGroup{}.Domain(group)), nil
}
