package mappers

import (
	domain "composition-api/internal/domain/exam"
	api "composition-api/internal/generated/http/api"
)

func Device(device domain.Device) api.ExamDevice {
	return api.ExamDevice{
		ID:   device.Id,
		Name: device.Name,
	}
}

func SliceDevice(devices []domain.Device) []api.ExamDevice {
	result := make([]api.ExamDevice, 0, len(devices))
	for _, device := range devices {
		result = append(result, Device(device))
	}
	return result
}
