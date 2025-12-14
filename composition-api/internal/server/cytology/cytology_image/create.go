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

func (h *handler) CytologyCreateCreate(ctx context.Context, req *api.CytologyCreateCreateReq) (api.CytologyCreateCreateRes, error) {
	arg := mappers.CytologyImage{}.CreateArgFromCytologyCreateCreateReq(req)

	id, err := h.services.CytologyService.CreateCytologyImage(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBadRequest):
			return &api.CytologyCreateCreateBadRequest{
				StatusCode: http.StatusBadRequest,
				Response: api.Error{
					Message: "Неверный формат запроса",
				},
			}, nil
		case errors.Is(err, domain.ErrUnprocessableEntity):
			return &api.CytologyCreateCreateUnprocessableEntity{
				StatusCode: http.StatusUnprocessableEntity,
				Response: api.Error{
					Message: "Ошибка валидации данных",
				},
			}, nil
		default:
			return nil, err
		}
	}

	// Возвращаем созданный объект согласно swagger.json
	result := api.CytologyCreateCreateCreated{
		DiagnosticNumber: req.DiagnosticNumber,
	}

	if req.ID.Set {
		result.ID = req.ID
	}
	if req.Image.Set {
		result.Image = req.Image
	}
	if req.IsLast.Set {
		result.IsLast = req.IsLast
	}
	if req.DiagnosDate.Set {
		result.DiagnosDate = req.DiagnosDate
	}
	if req.Details != nil {
		result.Details = &api.CytologyCreateCreateCreatedDetails{}
	}
	if req.DiagnosticMarking.Set {
		result.DiagnosticMarking = api.OptCytologyCreateCreateCreatedDiagnosticMarking{
			Value: api.CytologyCreateCreateCreatedDiagnosticMarking(req.DiagnosticMarking.Value),
			Set:   true,
		}
	}
	if req.MaterialType.Set {
		result.MaterialType = api.OptCytologyCreateCreateCreatedMaterialType{
			Value: api.CytologyCreateCreateCreatedMaterialType(req.MaterialType.Value),
			Set:   true,
		}
	}
	if req.Calcitonin.Set {
		result.Calcitonin = req.Calcitonin
	}
	if req.CalcitoninInFlush.Set {
		result.CalcitoninInFlush = req.CalcitoninInFlush
	}
	if req.Thyroglobulin.Set {
		result.Thyroglobulin = req.Thyroglobulin
	}
	if req.Prev.Set {
		result.Prev = req.Prev
	}
	if req.ParentPrev.Set {
		result.ParentPrev = req.ParentPrev
	}
	if req.PatientCard.Set {
		result.PatientCard = req.PatientCard
	}

	return pointer.To(result), nil
}
