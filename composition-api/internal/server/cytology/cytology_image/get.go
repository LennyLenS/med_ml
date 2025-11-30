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

	cytologyImage := mappers.CytologyImage{}.Domain(img)
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

	return pointer.To(api.CytologiesExternalIDGetOKApplicationJSON(
		mappers.CytologyImage{}.SliceDomain(imgs),
	)), nil
}

func (h *handler) CytologiesPatientCardDoctorIDPatientIDGet(ctx context.Context, params api.CytologiesPatientCardDoctorIDPatientIDGetParams) (api.CytologiesPatientCardDoctorIDPatientIDGetRes, error) {
	imgs, err := h.services.CytologyService.GetCytologyImagesByDoctorIdAndPatientId(ctx, params.DoctorID, params.PatientID)
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
		mappers.CytologyImage{}.SliceDomain(imgs),
	)), nil
}
