package mappers

import (
	domain "composition-api/internal/domain/exam"
	"composition-api/internal/generated/http/api"
)

type KT struct{}

func (KT) Domain(kt domain.KT) api.Kt {
	probabilities := kt.PredictedClasses
	if probabilities == nil {
		probabilities = make(map[string]float64)
	}
	return api.Kt{
		ID:                 kt.Id,
		CreateAt:           kt.CreateAt,
		Checked:            kt.Checked,
		AuthorID:           kt.Author,
		DeviceID:           kt.DeviceID,
		Status:             api.KtStatus(kt.Status),
		ClassProbabilities: api.KtClassProbabilities(probabilities),
	}
}

func (KT) SliceDomain(kts []domain.KT) []api.Kt {
	return slice(kts, KT{})
}
