package segmentation

import (
	"context"
	"errors"
	"net/http"

	"composition-api/internal/domain"

	api "composition-api/internal/generated/http/api"
)

func (h *handler) CytologySegmentationGroupIDDelete(ctx context.Context, params api.CytologySegmentationGroupIDDeleteParams) (api.CytologySegmentationGroupIDDeleteRes, error) {
	err := h.services.CytologyService.DeleteSegmentationGroup(ctx, params.ID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &api.CytologySegmentationGroupIDDeleteNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Группа сегментаций не найдена",
				},
			}, nil
		default:
			return nil, err
		}
	}

	return &api.CytologySegmentationGroupIDDeleteOK{}, nil
}
