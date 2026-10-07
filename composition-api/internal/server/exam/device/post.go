package device

import (
	"context"

	api "composition-api/internal/generated/http/api"
)

func (h *handler) MriDevicePost(ctx context.Context, req *api.MriDevicePostReq) (api.MriDevicePostRes, error) {
	id, err := h.services.ExamDeviceService.Create(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return &api.MriDevicePostOK{ID: id}, nil
}
