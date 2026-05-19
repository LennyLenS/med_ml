package mappers_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	cytology_domain "composition-api/internal/domain/cytology"
	med_domain "composition-api/internal/domain/med"
	api "composition-api/internal/generated/http/api"
	"composition-api/internal/server/cytology/mappers"
)

func TestSplitFullName_ThreeParts(t *testing.T) {
	patient := med_domain.Patient{FullName: "Иванов Иван Иванович"}
	result := mappers.CytologyImage{}.ToCytologyShotPatient(patient)

	require.Equal(t, "Иван", result.FirstName)
	require.Equal(t, "Иванов", result.LastName)
	require.Equal(t, "Иванович", result.FathersName)
}

func TestParseCytologyShotDetails(t *testing.T) {
	details := `{"ai_info":[],"probs":[0.1,0.5,0.4]}`
	marking := cytology_domain.DiagnosticMarkingP11
	material := cytology_domain.MaterialTypeTP
	calcitonin := 1
	imageID := uuid.New()
	originalID := uuid.New()
	diagnosis := "диагноз"

	shot := mappers.CytologyImage{}.ToCytologyPatientShot(
		cytology_domain.CytologyImage{
			Id:                imageID,
			DiagnosticNumber:  1,
			DiagnosticMarking: &marking,
			MaterialType:      &material,
			DiagnosDate:       time.Date(2026, 5, 17, 19, 51, 14, 0, time.UTC),
			IsLast:            true,
			Calcitonin:        &calcitonin,
			CalcitoninInFlush: &calcitonin,
			Thyroglobulin:     &calcitonin,
			Details:           func() *string { s := details; return &s }(),
		},
		med_domain.Card{Diagnosis: &diagnosis},
		&originalID,
	)

	require.Equal(t, imageID, shot.ID)
	require.True(t, shot.IsLast)
	require.Equal(t, []float64{0.1, 0.5, 0.4}, shot.Details.Probs)
	require.Empty(t, shot.Details.AiInfo)
	require.True(t, shot.DiagnosticMarking.Set)
	require.Equal(t, api.CytologyPatientShotDiagnosticMarking11, shot.DiagnosticMarking.Value)
	require.True(t, shot.OriginalImage.Set)
	require.False(t, shot.OriginalImage.Null)
	require.Equal(t, originalID, shot.OriginalImage.Value)
	require.True(t, shot.Prev.Null)
	require.True(t, shot.PatientCard.Diagnosis.Set)
	require.Equal(t, "диагноз", shot.PatientCard.Diagnosis.Value)
}
