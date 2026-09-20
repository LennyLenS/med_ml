package mappers

import (
	domain "composition-api/internal/domain/med"
	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/mappers"

	"github.com/google/uuid"
)

type Card struct{}

func (m Card) Api(d domain.Card) api.Card {
	card := api.Card{
		DoctorID:  d.DoctorID,
		PatientID: d.PatientID,
		Diagnosis: mappers.ToOptString(d.Diagnosis),
	}
	if d.UUID != uuid.Nil {
		card.UUID = api.OptUUID{
			Value: d.UUID,
			Set:   true,
		}
	}
	return card
}

func (m Card) SliceApi(d []domain.Card) []api.Card {
	return slice(d, m)
}
