package domain

import (
	"time"

	"github.com/google/uuid"
)

type Mri struct {
	Id          uuid.UUID
	Projection  MriProjection
	Checked     bool
	ExternalID  uuid.UUID
	Author      uuid.UUID
	DeviceID    int
	Status      ExamStatus
	Description *string
	CreateAt    time.Time
}
