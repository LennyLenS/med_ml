package mappers

import (
	domain "composition-api/internal/domain/exam"
	"composition-api/internal/generated/http/api"
)

type KT struct{}

func (KT) Domain(kt domain.KT) api.Kt {
	return api.Kt{
		ID:       kt.Id,
		CreateAt: kt.CreateAt,
	}
}

func (KT) SliceDomain(kts []domain.KT) []api.Kt {
	return slice(kts, KT{})
}
