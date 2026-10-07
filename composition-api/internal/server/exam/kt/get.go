package kt

import (
	"context"

	"composition-api/internal/generated/http/api"
	"composition-api/internal/server/exam/mappers"

	"github.com/AlekSi/pointer"
)

func (h *handler) KtIDGet(ctx context.Context, params api.KtIDGetParams) (api.KtIDGetRes, error) {
	kt, err := h.services.KTService.GetByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return pointer.To(mappers.KT{}.Domain(kt)), nil
}
