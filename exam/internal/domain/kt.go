package domain

import (
	"time"

	"github.com/google/uuid"
)

type KT struct {
	Id               uuid.UUID
	Checked          bool
	Author           uuid.UUID
	DeviceID         int
	Status           ExamStatus
	Description      *string
	CreateAt         time.Time
	PredictedClasses map[string]float64
}
