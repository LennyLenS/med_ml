package mappers

import (
	domain "composition-api/internal/domain/exam"
	api "composition-api/internal/generated/http/api"
)

type Mri struct{}

func (Mri) Domain(mri domain.Mri) api.Mri {
	return api.Mri{
		ID:         mri.Id,
		Projection: api.MriProjection(mri.Projection),
		Checked:    mri.Checked,
		ExternalID: mri.ExternalID,
		AuthorID:   mri.Author,
		DeviceID:   mri.DeviceID,
		Status:     api.MriStatus(mri.Status),
		CreateAt:   mri.CreateAt,
	}
}

func (Mri) SliceDomain(mris []domain.Mri) []api.Mri {
	return slice(mris, Mri{})
}
