package entity

import (
	"database/sql"

	gtclib "github.com/WantBeASleep/med_ml_lib/gtc"

	"med/internal/domain"

	"github.com/google/uuid"
)

type Card struct {
	ID        sql.NullInt32  `db:"id"`
	DoctorID  uuid.UUID      `db:"doctor_id"`
	PatientID uuid.UUID      `db:"patient_id"`
	Diagnosis sql.NullString `db:"diagnosis"`
}

func (Card) FromDomain(p domain.Card) Card {
	return Card{
		ID:        gtclib.Int.PointerToSql(p.ID),
		DoctorID:  p.DoctorID,
		PatientID: p.PatientID,
		Diagnosis: gtclib.String.PointerToSql(p.Diagnosis),
	}
}

func (p Card) ToDomain() domain.Card {
	return domain.Card{
		ID:        gtclib.Int.SqlToPointer(p.ID),
		DoctorID:  p.DoctorID,
		PatientID: p.PatientID,
		Diagnosis: gtclib.String.SqlToPointer(p.Diagnosis),
	}
}
