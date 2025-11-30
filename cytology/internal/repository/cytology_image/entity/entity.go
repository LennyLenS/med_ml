package entity

import (
	"database/sql"
	"encoding/json"
	"time"

	"cytology/internal/domain"

	"github.com/WantBeASleep/med_ml_lib/gtc"
	"github.com/google/uuid"
)

type CytologyImage struct {
	Id                uuid.UUID       `db:"id"`
	ExternalID        uuid.UUID       `db:"external_id"`
	PatientCardID     uuid.UUID       `db:"patient_card_id"`
	DiagnosticNumber  int             `db:"diagnostic_number"`
	DiagnosticMarking sql.NullString  `db:"diagnostic_marking"`
	MaterialType      sql.NullString  `db:"material_type"`
	DiagnosDate       time.Time       `db:"diagnos_date"`
	IsLast            bool            `db:"is_last"`
	Calcitonin        sql.NullInt32   `db:"calcitonin"`
	CalcitoninInFlush sql.NullInt32   `db:"calcitonin_in_flush"`
	Thyroglobulin     sql.NullInt32   `db:"thyroglobulin"`
	Details           json.RawMessage `db:"details"`
	PrevID            uuid.NullUUID   `db:"prev_id"`
	ParentPrevID      uuid.NullUUID   `db:"parent_prev_id"`
	CreateAt          time.Time       `db:"create_at"`
}

func (CytologyImage) FromDomain(d domain.CytologyImage) CytologyImage {
	var diagnosticMarking sql.NullString
	if d.DiagnosticMarking != nil {
		diagnosticMarking = sql.NullString{String: d.DiagnosticMarking.String(), Valid: true}
	}

	var materialType sql.NullString
	if d.MaterialType != nil {
		materialType = sql.NullString{String: d.MaterialType.String(), Valid: true}
	}

	return CytologyImage{
		Id:                d.Id,
		ExternalID:        d.ExternalID,
		PatientCardID:     d.PatientCardID,
		DiagnosticNumber:  d.DiagnosticNumber,
		DiagnosticMarking: diagnosticMarking,
		MaterialType:      materialType,
		DiagnosDate:       d.DiagnosDate,
		IsLast:            d.IsLast,
		Calcitonin:        gtc.Int32.PointerToSql(d.Calcitonin),
		CalcitoninInFlush: gtc.Int32.PointerToSql(d.CalcitoninInFlush),
		Thyroglobulin:     gtc.Int32.PointerToSql(d.Thyroglobulin),
		Details:           d.Details,
		PrevID:            gtc.UUID.PointerToNullUUID(d.PrevID),
		ParentPrevID:      gtc.UUID.PointerToNullUUID(d.ParentPrevID),
		CreateAt:          d.CreateAt,
	}
}

func (d CytologyImage) ToDomain() domain.CytologyImage {
	var diagnosticMarking *domain.DiagnosticMarking
	if d.DiagnosticMarking.Valid {
		dm := domain.DiagnosticMarking(d.DiagnosticMarking.String)
		diagnosticMarking = &dm
	}

	var materialType *domain.MaterialType
	if d.MaterialType.Valid {
		mt := domain.MaterialType(d.MaterialType.String)
		materialType = &mt
	}

	return domain.CytologyImage{
		Id:                d.Id,
		ExternalID:        d.ExternalID,
		PatientCardID:     d.PatientCardID,
		DiagnosticNumber:  d.DiagnosticNumber,
		DiagnosticMarking: diagnosticMarking,
		MaterialType:      materialType,
		DiagnosDate:       d.DiagnosDate,
		IsLast:            d.IsLast,
		Calcitonin:        gtc.Int32.SqlToPointer(d.Calcitonin),
		CalcitoninInFlush: gtc.Int32.SqlToPointer(d.CalcitoninInFlush),
		Thyroglobulin:     gtc.Int32.SqlToPointer(d.Thyroglobulin),
		Details:           d.Details,
		PrevID:            gtc.UUID.NullUUIDToPointer(d.PrevID),
		ParentPrevID:      gtc.UUID.NullUUIDToPointer(d.ParentPrevID),
		CreateAt:          d.CreateAt,
	}
}

func (CytologyImage) SliceToDomain(images []CytologyImage) []domain.CytologyImage {
	domainImages := make([]domain.CytologyImage, 0, len(images))
	for _, v := range images {
		domainImages = append(domainImages, v.ToDomain())
	}
	return domainImages
}
