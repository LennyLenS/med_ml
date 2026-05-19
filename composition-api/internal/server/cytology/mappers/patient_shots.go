package mappers

import (
	"encoding/json"
	"strings"

	"github.com/go-faster/jx"
	"github.com/google/uuid"

	cytology_domain "composition-api/internal/domain/cytology"
	med_domain "composition-api/internal/domain/med"
	api "composition-api/internal/generated/http/api"
)

type cytologyShotDetailsJSON struct {
	AIInfo []json.RawMessage `json:"ai_info"`
	Probs  []float64         `json:"probs"`
}

func splitFullName(fullname string) (firstName, lastName, fathersName string) {
	parts := strings.Fields(strings.TrimSpace(fullname))
	switch len(parts) {
	case 0:
		return "", "", ""
	case 1:
		return parts[0], "", ""
	case 2:
		return parts[1], parts[0], ""
	default:
		return parts[1], parts[0], strings.Join(parts[2:], " ")
	}
}

func parseCytologyShotDetails(details *string) api.CytologyShotDetails {
	result := api.CytologyShotDetails{
		AiInfo: []jx.Raw{},
		Probs:  []float64{},
	}
	if details == nil || *details == "" {
		return result
	}

	var parsed cytologyShotDetailsJSON
	if err := json.Unmarshal([]byte(*details), &parsed); err != nil {
		return result
	}

	if len(parsed.AIInfo) > 0 {
		result.AiInfo = make([]jx.Raw, 0, len(parsed.AIInfo))
		for _, item := range parsed.AIInfo {
			result.AiInfo = append(result.AiInfo, jx.Raw(item))
		}
	}

	if len(parsed.Probs) > 0 {
		result.Probs = parsed.Probs
	}

	return result
}

func (CytologyImage) ToCytologyShotPatient(patient med_domain.Patient) api.CytologyShotPatient {
	firstName, lastName, fathersName := splitFullName(patient.FullName)

	return api.CytologyShotPatient{
		ID:             patient.Id,
		FirstName:      firstName,
		LastName:       lastName,
		FathersName:    fathersName,
		BirthDate:      patient.BirthDate,
		PersonalPolicy: patient.Policy,
		Email:          patient.Email,
		IsActive:       patient.Active,
	}
}

func (CytologyImage) ToCytologyPatientShot(
	img cytology_domain.CytologyImage,
	patientCard med_domain.Card,
	originalImageID *uuid.UUID,
) api.CytologyPatientShot {
	shot := api.CytologyPatientShot{
		ID:               img.Id,
		PatientCard:      toCytologyShotPatientCard(patientCard),
		IsLast:           img.IsLast,
		DiagnosDate:      img.DiagnosDate,
		Details:          parseCytologyShotDetails(img.Details),
		DiagnosticNumber: img.DiagnosticNumber,
	}

	if img.DiagnosticMarking != nil {
		shot.DiagnosticMarking = api.OptCytologyPatientShotDiagnosticMarking{
			Value: api.CytologyPatientShotDiagnosticMarking(*img.DiagnosticMarking),
			Set:   true,
		}
	}

	if img.MaterialType != nil {
		shot.MaterialType = api.OptCytologyPatientShotMaterialType{
			Value: api.CytologyPatientShotMaterialType(*img.MaterialType),
			Set:   true,
		}
	}

	if img.Calcitonin != nil {
		shot.Calcitonin = api.OptInt{
			Value: int(*img.Calcitonin),
			Set:   true,
		}
	}

	if img.CalcitoninInFlush != nil {
		shot.CalcitoninInFlush = api.OptInt{
			Value: int(*img.CalcitoninInFlush),
			Set:   true,
		}
	}

	if img.Thyroglobulin != nil {
		shot.Thyroglobulin = api.OptInt{
			Value: int(*img.Thyroglobulin),
			Set:   true,
		}
	}

	if img.PrevID != nil {
		shot.Prev.SetTo(*img.PrevID)
	} else {
		shot.Prev.SetToNull()
	}

	if img.ParentPrevID != nil {
		shot.ParentPrev.SetTo(*img.ParentPrevID)
	} else {
		shot.ParentPrev.SetToNull()
	}

	if originalImageID != nil {
		shot.OriginalImage.SetTo(*originalImageID)
	} else {
		shot.OriginalImage.SetToNull()
	}

	return shot
}

func toCytologyShotPatientCard(card med_domain.Card) api.CytologyShotPatientCard {
	result := api.CytologyShotPatientCard{}

	if card.Diagnosis != nil {
		result.Diagnosis = api.OptString{
			Value: *card.Diagnosis,
			Set:   true,
		}
	}

	return result
}
