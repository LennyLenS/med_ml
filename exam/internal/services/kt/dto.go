package kt

import (
	"github.com/google/uuid"

	"exam/internal/domain"
)

type CreateKTArg struct {
	Author      uuid.UUID
	DeviceID    int
	Description *string
}

type UpdateKTArg struct {
	Id               uuid.UUID
	Checked          *bool
	PredictedClasses *map[string]float64
}

func (u UpdateKTArg) UpdateDomain(d *domain.KT) {
	if u.Checked != nil {
		d.Checked = *u.Checked
	}
	if u.PredictedClasses != nil {
		d.PredictedClasses = *u.PredictedClasses
	}
}
