package cytology_image

import (
	"context"
	"errors"
	"net/http"

	"composition-api/internal/domain"

	"github.com/AlekSi/pointer"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/cytology/mappers"
)

func (h *handler) CytologyIDGet(ctx context.Context, params api.CytologyIDGetParams) (api.CytologyIDGetRes, error) {
	img, origImg, err := h.services.CytologyService.GetCytologyImageById(ctx, params.ID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &api.CytologyIDGetNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Цитологическое исследование не найдено",
				},
			}, nil
		default:
			return nil, err
		}
	}

	// TODO: Get doctor_id and patient_id from patient_card_id
	// For now, we'll use zero UUIDs as placeholders
	// This needs to be fixed by either:
	// 1. Storing doctor_id and patient_id in cytology microservice
	// 2. Using med service to get doctor_id and patient_id from patient_card_id
	var doctorID, patientID = img.Id, img.Id // Placeholder

	cytologyImage := mappers.CytologyImage{}.Domain(img, doctorID, patientID)
	result := &api.CytologyIDGetOK{
		CytologyImage: api.OptCytologyImage{
			Value: cytologyImage,
			Set:   true,
		},
	}

	if origImg != nil {
		originalImage := mappers.OriginalImage{}.Domain(*origImg)
		result.OriginalImage = api.OptOriginalImage{
			Value: originalImage,
			Set:   true,
		}
	}

	return result, nil
}

func (h *handler) CytologiesExternalIDGet(ctx context.Context, params api.CytologiesExternalIDGetParams) (api.CytologiesExternalIDGetRes, error) {
	imgs, err := h.services.CytologyService.GetCytologyImagesByExternalId(ctx, params.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &api.CytologiesExternalIDGetNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Цитологические исследования не найдены",
				},
			}, nil
		}
		return nil, err
	}

	// TODO: Get doctor_id and patient_id from patient_card_id for each image
	// For now, we'll use zero UUIDs as placeholders
	var doctorID, patientID = params.ID, params.ID // Placeholder

	return pointer.To(api.CytologiesExternalIDGetOKApplicationJSON(
		mappers.CytologyImage{}.SliceDomain(imgs, doctorID, patientID),
	)), nil
}

func (h *handler) CytologiesPatientCardDoctorIDPatientIDGet(ctx context.Context, params api.CytologiesPatientCardDoctorIDPatientIDGetParams) (api.CytologiesPatientCardDoctorIDPatientIDGetRes, error) {
	patientCardID := mappers.GetPatientCardID(params.DoctorID, params.PatientID)

	imgs, err := h.services.CytologyService.GetCytologyImagesByPatientCardId(ctx, patientCardID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &api.CytologiesPatientCardDoctorIDPatientIDGetNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Цитологические исследования не найдены",
				},
			}, nil
		}
		return nil, err
	}

	return pointer.To(api.CytologiesPatientCardDoctorIDPatientIDGetOKApplicationJSON(
		mappers.CytologyImage{}.SliceDomain(imgs, params.DoctorID, params.PatientID),
	)), nil
}
