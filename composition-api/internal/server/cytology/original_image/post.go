package original_image

import (
	"context"
	"errors"
	"net/http"

	"composition-api/internal/domain"

	"github.com/AlekSi/pointer"

	api "composition-api/internal/generated/http/api"
	cytologySrv "composition-api/internal/services/cytology"
)

func (h *handler) CytologyIDOriginalImagePost(ctx context.Context, req *api.CytologyIDOriginalImagePostReq, params api.CytologyIDOriginalImagePostParams) (api.CytologyIDOriginalImagePostRes, error) {
	contentType := req.File.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	arg := cytologySrv.CreateOriginalImageArg{
		CytologyID:  params.ID,
		File:        req.File,
		ContentType: contentType,
	}

	if req.DelayTime.Set {
		arg.DelayTime = &req.DelayTime.Value
	}

	id, err := h.services.CytologyService.CreateOriginalImage(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadRequest):
			return &api.CytologyIDOriginalImagePostBadRequest{
				StatusCode: http.StatusBadRequest,
				Response: api.Error{
					Message: "Неверный формат запроса",
				},
			}, nil
		case errors.Is(err, domain.ErrNotFound):
			return &api.CytologyIDOriginalImagePostNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Цитологическое исследование не найдено",
				},
			}, nil
		default:
			return nil, err
		}
	}

	return pointer.To(api.SimpleUuid{ID: id}), nil
}
