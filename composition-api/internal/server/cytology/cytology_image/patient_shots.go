package cytology_image

import (
	"context"
	"errors"
	"net/http"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"

	"composition-api/internal/domain"
	auth_domain "composition-api/internal/domain/auth"
	cytology_domain "composition-api/internal/domain/cytology"
	med_domain "composition-api/internal/domain/med"
	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/cytology/mappers"
	"composition-api/internal/server/security"
)

func (h *handler) CytologyPatientShotsRead(ctx context.Context, params api.CytologyPatientShotsReadParams) (api.CytologyPatientShotsReadRes, error) {
	patientID := params.PatientID

	token, err := security.ParseToken(ctx)
	if err != nil {
		return nil, err
	}

	var images []cytology_domain.CytologyImage
	switch token.Role {
	case auth_domain.RoleDoctor:
		images, err = h.services.CytologyService.GetCytologyImagesByDoctorIdAndPatientId(ctx, token.Id, patientID)
	case auth_domain.RolePatient:
		if token.Id != patientID {
			return &api.CytologyPatientShotsReadForbidden{
				StatusCode: http.StatusForbidden,
				Response: api.Error{
					Message: "Доступ запрещён",
				},
			}, nil
		}
		images, err = h.services.CytologyService.GetCytologyImagesByPatientId(ctx, patientID)
	default:
		return &api.CytologyPatientShotsReadForbidden{
			StatusCode: http.StatusForbidden,
			Response: api.Error{
				Message: "Доступ запрещён",
			},
		}, nil
	}
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	patient, err := h.services.PatientService.GetPatient(ctx, patientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &api.CytologyPatientShotsReadNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Пациент не найден",
				},
			}, nil
		}
		return nil, err
	}

	shots, err := h.buildPatientShots(ctx, images)
	if err != nil {
		return nil, err
	}

	return pointer.To(api.CytologyPatientShotsReadOK{
		Patient: mappers.CytologyImage{}.ToCytologyShotPatient(patient),
		Shots:   shots,
	}), nil
}

func (h *handler) buildPatientShots(ctx context.Context, images []cytology_domain.CytologyImage) ([]api.CytologyPatientShot, error) {
	shots := make([]api.CytologyPatientShot, 0, len(images))

	for _, img := range images {
		var patientCard med_domain.Card
		if img.DoctorID != uuid.Nil && img.PatientID != uuid.Nil {
			card, err := h.services.CardService.GetCard(ctx, img.DoctorID, img.PatientID)
			if err == nil {
				patientCard = card
			}
		}

		var originalImageID *uuid.UUID
		originalImages, err := h.services.CytologyService.GetOriginalImagesByCytologyId(ctx, img.Id)
		if err == nil && len(originalImages) > 0 {
			originalImageID = &originalImages[0].Id
		}

		shots = append(shots, mappers.CytologyImage{}.ToCytologyPatientShot(img, patientCard, originalImageID))
	}

	return shots, nil
}
