package domain

import "github.com/google/uuid"

type Card struct {
	ID        *int
	UUID      uuid.UUID
	DoctorID  uuid.UUID
	PatientID uuid.UUID
	Diagnosis *string
}

func (c Card) Validate() error {
	if c.DoctorID == uuid.Nil || c.PatientID == uuid.Nil {
		return ErrBadRequest
	}
	return nil
}
