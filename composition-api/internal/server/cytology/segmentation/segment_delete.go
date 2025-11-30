package segmentation

import (
	"context"
	"errors"
	"net/http"

	"composition-api/internal/domain"

	api "composition-api/internal/generated/http/api"
)

func (h *handler) CytologySegmentationIDDelete(ctx context.Context, params api.CytologySegmentationIDDeleteParams) (api.CytologySegmentationIDDeleteRes, error) {
	err := h.services.CytologyService.DeleteSegmentation(ctx, params.ID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &api.CytologySegmentationIDDeleteNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Сегмент не найден",
				},
			}, nil
		default:
			return nil, err
		}
	}

	return &api.CytologySegmentationIDDeleteOK{}, nil
}
