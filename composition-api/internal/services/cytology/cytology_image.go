package cytology

import (
	"context"

	"github.com/google/uuid"

	"composition-api/internal/adapters/cytology"
	domain "composition-api/internal/domain/cytology"
)

func (s *service) CreateCytologyImage(ctx context.Context, arg CreateCytologyImageArg) (uuid.UUID, error) {
	return s.adapters.Cytology.CreateCytologyImage(ctx, cytology.CreateCytologyImageIn{
		ExternalID:        arg.ExternalID,
		DoctorID:          arg.DoctorID,
		PatientID:         arg.PatientID,
		DiagnosticNumber:  arg.DiagnosticNumber,
		DiagnosticMarking: arg.DiagnosticMarking,
		MaterialType:      arg.MaterialType,
		Calcitonin:        arg.Calcitonin,
		CalcitoninInFlush: arg.CalcitoninInFlush,
		Thyroglobulin:     arg.Thyroglobulin,
		Details:           arg.Details,
		PrevID:            arg.PrevID,
		ParentPrevID:      arg.ParentPrevID,
	})
}

func (s *service) GetCytologyImageById(ctx context.Context, id uuid.UUID) (domain.CytologyImage, *domain.OriginalImage, error) {
	return s.adapters.Cytology.GetCytologyImageById(ctx, id)
}

func (s *service) GetCytologyImagesByExternalId(ctx context.Context, externalID uuid.UUID) ([]domain.CytologyImage, error) {
	return s.adapters.Cytology.GetCytologyImagesByExternalId(ctx, externalID)
}

func (s *service) GetCytologyImagesByDoctorIdAndPatientId(ctx context.Context, doctorID, patientID uuid.UUID) ([]domain.CytologyImage, error) {
	return s.adapters.Cytology.GetCytologyImagesByDoctorIdAndPatientId(ctx, doctorID, patientID)
}

func (s *service) UpdateCytologyImage(ctx context.Context, arg UpdateCytologyImageArg) (domain.CytologyImage, error) {
	return s.adapters.Cytology.UpdateCytologyImage(ctx, cytology.UpdateCytologyImageIn{
		Id:                arg.Id,
		DiagnosticMarking: arg.DiagnosticMarking,
		MaterialType:      arg.MaterialType,
		Calcitonin:        arg.Calcitonin,
		CalcitoninInFlush: arg.CalcitoninInFlush,
		Thyroglobulin:     arg.Thyroglobulin,
		Details:           arg.Details,
		IsLast:            arg.IsLast,
	})
}

func (s *service) DeleteCytologyImage(ctx context.Context, id uuid.UUID) error {
	return s.adapters.Cytology.DeleteCytologyImage(ctx, id)
}
