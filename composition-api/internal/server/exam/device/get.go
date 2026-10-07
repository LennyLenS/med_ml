package device

import (
	"context"

	"github.com/AlekSi/pointer"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/exam/mappers"
)

func (h *handler) MriDevicesGet(ctx context.Context) (api.MriDevicesGetRes, error) {
	devices, err := h.services.ExamDeviceService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return pointer.To(api.MriDevicesGetOKApplicationJSON(mappers.SliceDevice(devices))), nil
}
