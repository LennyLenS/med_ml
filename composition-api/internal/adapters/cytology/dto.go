package cytology

import (
	"github.com/google/uuid"

	domain "composition-api/internal/domain/cytology"
)

type CreateCytologyImageIn struct {
	ExternalID        uuid.UUID
	PatientCardID     uuid.UUID
	DiagnosticNumber  int
	DiagnosticMarking *domain.DiagnosticMarking
	MaterialType      *domain.MaterialType
	Calcitonin        *int
	CalcitoninInFlush *int
	Thyroglobulin     *int
	Details           *string
	PrevID            *uuid.UUID
	ParentPrevID      *uuid.UUID
}

type UpdateCytologyImageIn struct {
	Id                uuid.UUID
	DiagnosticMarking *domain.DiagnosticMarking
	MaterialType      *domain.MaterialType
	Calcitonin        *int
	CalcitoninInFlush *int
	Thyroglobulin     *int
	Details           *string
	IsLast            *bool
}

type CreateOriginalImageIn struct {
	CytologyID uuid.UUID
	ImagePath  string
	DelayTime  *float64
}

type UpdateOriginalImageIn struct {
	Id         uuid.UUID
	DelayTime  *float64
	ViewedFlag *bool
}

type CreateSegmentationGroupIn struct {
	CytologyID uuid.UUID
	SegType    domain.SegType
	GroupType  domain.GroupType
	IsAI       bool
	Details    *string
}

type UpdateSegmentationGroupIn struct {
	Id      uuid.UUID
	SegType *domain.SegType
	Details *string
}

type CreateSegmentationIn struct {
	SegmentationGroupID uuid.UUID
	Points              []domain.SegmentationPoint
}

type UpdateSegmentationIn struct {
	Id     uuid.UUID
	Points []domain.SegmentationPoint
}
