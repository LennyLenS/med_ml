package kt

import (
	"context"

	"github.com/AlekSi/pointer"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/exam/mappers"
	apimappers "composition-api/internal/server/mappers"
	ktSrv "composition-api/internal/services/kt"
)

func (h *handler) KtIDPatch(ctx context.Context, req *api.KtIDPatchReq, params api.KtIDPatchParams) (api.KtIDPatchRes, error) {
	kt, err := h.services.KTService.Update(ctx, ktSrv.UpdateKtArg{
		Id:      params.ID,
		Checked: apimappers.FromOptBool(req.Checked),
	})
	if err != nil {
		return nil, err
	}
	return pointer.To(mappers.KT{}.Domain(kt)), nil
}
