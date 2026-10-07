package mri

import (
	"context"

	api "composition-api/internal/generated/http/api"
)

func (h *handler) MriIDDelete(ctx context.Context, params api.MriIDDeleteParams) (api.MriIDDeleteRes, error) {
	err := h.services.MriService.DeleteByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.MriIDDeleteOK{}, nil
}
