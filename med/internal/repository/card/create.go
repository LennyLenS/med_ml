package card

import (
	"med/internal/repository/card/entity"
	repoEntity "med/internal/repository/entity"

	"github.com/google/uuid"
)

func (r *repo) InsertCard(card entity.Card) (int, uuid.UUID, error) {
	query := r.QueryBuilder().
		Insert(table).
		Columns(
			columnDoctorID,
			columnPatientID,
			columnDiagnosis,
		).
		Values(
			card.DoctorID,
			card.PatientID,
			card.Diagnosis,
		).
		Suffix("RETURNING id, uuid")

	var res struct {
		ID   int       `db:"id"`
		UUID uuid.UUID `db:"uuid"`
	}
	err := r.Runner().Getx(r.Context(), &res, query)
	if err != nil {
		return 0, uuid.Nil, repoEntity.WrapDBError(err)
	}

	return res.ID, res.UUID, nil
}
