package original_image

import (
	"context"
	"errors"
	"net/http"

	"composition-api/internal/domain"

	"github.com/AlekSi/pointer"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/cytology/mappers"
	cytologySrv "composition-api/internal/services/cytology"
)

func (h *handler) CytologyOriginalImageIDPatch(ctx context.Context, req *api.CytologyOriginalImageIDPatchReq, params api.CytologyOriginalImageIDPatchParams) (api.CytologyOriginalImageIDPatchRes, error) {
	arg := cytologySrv.UpdateOriginalImageArg{
		Id: params.ID,
	}

	if req.DelayTime.Set {
		arg.DelayTime = &req.DelayTime.Value
	}

	if req.ViewedFlag.Set {
		arg.ViewedFlag = &req.ViewedFlag.Value
	}

	img, err := h.services.CytologyService.UpdateOriginalImage(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &api.CytologyOriginalImageIDPatchNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Оригинальное изображение не найдено",
				},
			}, nil
		default:
			return nil, err
		}
	}

	return pointer.To(mappers.OriginalImage{}.Domain(img)), nil
}
