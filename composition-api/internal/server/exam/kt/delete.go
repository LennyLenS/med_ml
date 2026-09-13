package kt

import (
	"context"

	api "composition-api/internal/generated/http/api"
)

func (h *handler) KtIDDelete(ctx context.Context, params api.KtIDDeleteParams) (api.KtIDDeleteRes, error) {
	err := h.services.KTService.DeleteByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.KtIDDeleteOK{}, nil
}
