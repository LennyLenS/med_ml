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

func (h *handler) CytologyIDPatch(ctx context.Context, req *api.CytologyIDPatchReq, params api.CytologyIDPatchParams) (api.CytologyIDPatchRes, error) {
	arg := mappers.CytologyImage{}.UpdateArg(params.ID, req)

	img, err := h.services.CytologyService.UpdateCytologyImage(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return &api.CytologyIDPatchNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Цитологическое исследование не найдено",
				},
			}, nil
		case errors.Is(err, domain.ErrBadRequest):
			return &api.CytologyIDPatchBadRequest{
				StatusCode: http.StatusBadRequest,
				Response: api.Error{
					Message: "Неверный формат запроса",
				},
			}, nil
		default:
			return nil, err
		}
	}

	// TODO: Get doctor_id and patient_id from patient_card_id
	// For now, we'll use zero UUIDs as placeholders
	var doctorID, patientID = img.Id, img.Id // Placeholder

	result := mappers.CytologyImage{}.Domain(img, doctorID, patientID)
	return pointer.To(result), nil
}
