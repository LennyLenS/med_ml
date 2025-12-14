package mappers

import (
	"github.com/google/uuid"

	domain "composition-api/internal/domain/cytology"
	api "composition-api/internal/generated/http/api"
	cytologySrv "composition-api/internal/services/cytology"
)

type CytologyImage struct{}

func (CytologyImage) Domain(img domain.CytologyImage) api.CytologyImage {
	result := api.CytologyImage{
		ID:               img.Id,
		ExternalID:       img.ExternalID,
		DoctorID:         img.DoctorID,
		PatientID:        img.PatientID,
		DiagnosticNumber: int(img.DiagnosticNumber),
		DiagnosDate:      img.DiagnosDate,
		IsLast:           img.IsLast,
		CreateAt:         img.CreateAt,
	}

	if img.DiagnosticMarking != nil {
		result.DiagnosticMarking = api.OptCytologyImageDiagnosticMarking{
			Value: api.CytologyImageDiagnosticMarking(*img.DiagnosticMarking),
			Set:   true,
		}
	}

	if img.MaterialType != nil {
		result.MaterialType = api.OptCytologyImageMaterialType{
			Value: api.CytologyImageMaterialType(*img.MaterialType),
			Set:   true,
		}
	}

	if img.Calcitonin != nil {
		result.Calcitonin = api.OptInt{
			Value: int(*img.Calcitonin),
			Set:   true,
		}
	}

	if img.CalcitoninInFlush != nil {
		result.CalcitoninInFlush = api.OptInt{
			Value: int(*img.CalcitoninInFlush),
			Set:   true,
		}
	}

	if img.Thyroglobulin != nil {
		result.Thyroglobulin = api.OptInt{
			Value: int(*img.Thyroglobulin),
			Set:   true,
		}
	}

	return result
}

func (CytologyImage) SliceDomain(imgs []domain.CytologyImage) []api.CytologyImage {
	result := make([]api.CytologyImage, 0, len(imgs))
	for _, img := range imgs {
		result = append(result, CytologyImage{}.Domain(img))
	}
	return result
}

func (CytologyImage) CreateArg(req *api.CytologyPostReq) cytologySrv.CreateCytologyImageArg {
	arg := cytologySrv.CreateCytologyImageArg{
		ExternalID:       req.ExternalID,
		DoctorID:         req.DoctorID,
		PatientID:        req.PatientID,
		DiagnosticNumber: int(req.DiagnosticNumber),
	}

	if req.DiagnosticMarking.Set {
		marking := domain.DiagnosticMarking(req.DiagnosticMarking.Value)
		arg.DiagnosticMarking = &marking
	}

	if req.MaterialType.Set {
		materialType := domain.MaterialType(req.MaterialType.Value)
		arg.MaterialType = &materialType
	}

	if req.Calcitonin.Set {
		calcitonin := int(req.Calcitonin.Value)
		arg.Calcitonin = &calcitonin
	}

	if req.CalcitoninInFlush.Set {
		calcitoninInFlush := int(req.CalcitoninInFlush.Value)
		arg.CalcitoninInFlush = &calcitoninInFlush
	}

	if req.Thyroglobulin.Set {
		thyroglobulin := int(req.Thyroglobulin.Value)
		arg.Thyroglobulin = &thyroglobulin
	}

	if req.Details != nil {
		// Details is a JSON string, we'll need to marshal it
		// For now, we'll pass it as is if it's already a string
		// This might need adjustment based on the actual structure
	}

	if req.PrevID.Set {
		arg.PrevID = &req.PrevID.Value
	}

	if req.ParentPrevID.Set {
		arg.ParentPrevID = &req.ParentPrevID.Value
	}

	return arg
}

func (CytologyImage) CreateArgFromCytologyCreateCreateReq(req *api.CytologyCreateCreateReq) cytologySrv.CreateCytologyImageArg {
	arg := cytologySrv.CreateCytologyImageArg{
		DiagnosticNumber: req.DiagnosticNumber,
		// ExternalID, DoctorID, PatientID должны быть получены из контекста или другого источника
		// Пока используем пустые UUID, но это нужно будет исправить
		ExternalID: uuid.Nil,
		DoctorID:   uuid.Nil,
		PatientID:  uuid.Nil,
	}

	if req.DiagnosticMarking.Set {
		marking := domain.DiagnosticMarking(req.DiagnosticMarking.Value)
		arg.DiagnosticMarking = &marking
	}

	if req.MaterialType.Set {
		materialType := domain.MaterialType(req.MaterialType.Value)
		arg.MaterialType = &materialType
	}

	if req.Calcitonin.Set {
		calcitonin := int(req.Calcitonin.Value)
		arg.Calcitonin = &calcitonin
	}

	if req.CalcitoninInFlush.Set {
		calcitoninInFlush := int(req.CalcitoninInFlush.Value)
		arg.CalcitoninInFlush = &calcitoninInFlush
	}

	if req.Thyroglobulin.Set {
		thyroglobulin := int(req.Thyroglobulin.Value)
		arg.Thyroglobulin = &thyroglobulin
	}

	if req.Details != nil {
		// Details handling - может быть JSON объект
	}

	// Prev и ParentPrev в swagger.json - это integer, но в нашей системе это UUID
	// Нужно будет преобразовать или получить из другого источника
	// Пока оставляем nil

	return arg
}

func (CytologyImage) UpdateArg(id uuid.UUID, req *api.CytologyIDUpdatePatchReq) cytologySrv.UpdateCytologyImageArg {
	arg := cytologySrv.UpdateCytologyImageArg{
		Id: id,
	}

	if req.DiagnosticMarking.Set {
		marking := domain.DiagnosticMarking(req.DiagnosticMarking.Value)
		arg.DiagnosticMarking = &marking
	}

	if req.MaterialType.Set {
		materialType := domain.MaterialType(req.MaterialType.Value)
		arg.MaterialType = &materialType
	}

	if req.Calcitonin.Set {
		calcitonin := int(req.Calcitonin.Value)
		arg.Calcitonin = &calcitonin
	}

	if req.CalcitoninInFlush.Set {
		calcitoninInFlush := int(req.CalcitoninInFlush.Value)
		arg.CalcitoninInFlush = &calcitoninInFlush
	}

	if req.Thyroglobulin.Set {
		thyroglobulin := int(req.Thyroglobulin.Value)
		arg.Thyroglobulin = &thyroglobulin
	}

	if req.Details != nil {
		// Details handling - might need to marshal to JSON string
	}

	if req.IsLast.Set {
		arg.IsLast = &req.IsLast.Value
	}

	return arg
}

func (CytologyImage) ToCytologyReadOKInfo(img domain.CytologyImage) api.CytologyReadOKInfo {
	// TODO: Patient и PatientCard должны быть получены из другого источника
	// Пока создаем заглушки
	patient := api.Patient{
		ID:       img.PatientID,
		Fullname: "", // Нужно получить из другого источника
		Email:    "", // Нужно получить из другого источника
		Policy:   "", // Нужно получить из другого источника
		Active:   true,
	}

	patientCard := api.PatientCard{
		Patient: api.OptInt{
			// Нужно преобразовать UUID в int или получить из другого источника
			Set: false,
		},
		MedWorker: api.OptInt{
			// Нужно преобразовать UUID в int или получить из другого источника
			Set: false,
		},
		Diagnosis: api.OptString{
			Set: false,
		},
	}

	imageGroup := api.CytologyReadOKInfoImageGroup{
		DiagnosticNumber: int(img.DiagnosticNumber),
	}

	if img.DiagnosticMarking != nil {
		imageGroup.DiagnosticMarking = api.OptCytologyReadOKInfoImageGroupDiagnosticMarking{
			Value: api.CytologyReadOKInfoImageGroupDiagnosticMarking(*img.DiagnosticMarking),
			Set:   true,
		}
	}

	if img.MaterialType != nil {
		imageGroup.MaterialType = api.OptCytologyReadOKInfoImageGroupMaterialType{
			Value: api.CytologyReadOKInfoImageGroupMaterialType(*img.MaterialType),
			Set:   true,
		}
	}

	if img.Calcitonin != nil {
		imageGroup.Calcitonin = api.OptInt{
			Value: int(*img.Calcitonin),
			Set:   true,
		}
	}

	if img.CalcitoninInFlush != nil {
		imageGroup.CalcitoninInFlush = api.OptInt{
			Value: int(*img.CalcitoninInFlush),
			Set:   true,
		}
	}

	if img.Thyroglobulin != nil {
		imageGroup.Thyroglobulin = api.OptInt{
			Value: int(*img.Thyroglobulin),
			Set:   true,
		}
	}

	imageGroup.IsLast = api.OptBool{
		Value: img.IsLast,
		Set:   true,
	}

	imageGroup.DiagnosDate = api.OptDateTime{
		Value: img.DiagnosDate,
		Set:   true,
	}

	return api.CytologyReadOKInfo{
		Patient:     patient,
		PatientCard: patientCard,
		ImageGroup:  imageGroup,
	}
}

func (CytologyImage) ToCytologyImageModelList(imgs []domain.CytologyImage) []api.CytologyHistoryReadOKResultsItem {
	result := make([]api.CytologyHistoryReadOKResultsItem, 0, len(imgs))
	for _, img := range imgs {
		item := api.CytologyHistoryReadOKResultsItem{
			DiagnosticNumber: int(img.DiagnosticNumber),
			IsLast: api.OptBool{
				Value: img.IsLast,
				Set:   true,
			},
			DiagnosDate: api.OptDateTime{
				Value: img.DiagnosDate,
				Set:   true,
			},
		}

		if img.DiagnosticMarking != nil {
			item.DiagnosticMarking = api.OptCytologyHistoryReadOKResultsItemDiagnosticMarking{
				Value: api.CytologyHistoryReadOKResultsItemDiagnosticMarking(*img.DiagnosticMarking),
				Set:   true,
			}
		}

		if img.MaterialType != nil {
			item.MaterialType = api.OptCytologyHistoryReadOKResultsItemMaterialType{
				Value: api.CytologyHistoryReadOKResultsItemMaterialType(*img.MaterialType),
				Set:   true,
			}
		}

		if img.Calcitonin != nil {
			item.Calcitonin = api.OptInt{
				Value: int(*img.Calcitonin),
				Set:   true,
			}
		}

		if img.CalcitoninInFlush != nil {
			item.CalcitoninInFlush = api.OptInt{
				Value: int(*img.CalcitoninInFlush),
				Set:   true,
			}
		}

		if img.Thyroglobulin != nil {
			item.Thyroglobulin = api.OptInt{
				Value: int(*img.Thyroglobulin),
				Set:   true,
			}
		}

		if img.Details != nil {
			item.Details = &api.CytologyHistoryReadOKResultsItemDetails{}
		}

		result = append(result, item)
	}
	return result
}

func (CytologyImage) UpdateArgFromCytologyUpdateUpdateReq(id uuid.UUID, req *api.CytologyUpdateUpdateReq) cytologySrv.UpdateCytologyImageArg {
	arg := cytologySrv.UpdateCytologyImageArg{
		Id: id,
	}

	// Извлекаем данные из req.Details или из верхнего уровня
	if req.Details.Set {
		details := req.Details.Value
		if details.DiagnosticMarking.Set {
			marking := domain.DiagnosticMarking(details.DiagnosticMarking.Value)
			arg.DiagnosticMarking = &marking
		}
		if details.MaterialType.Set {
			materialType := domain.MaterialType(details.MaterialType.Value)
			arg.MaterialType = &materialType
		}
		if details.Calcitonin.Set {
			calcitonin := int(details.Calcitonin.Value)
			arg.Calcitonin = &calcitonin
		}
		if details.CalcitoninInFlush.Set {
			calcitoninInFlush := int(details.CalcitoninInFlush.Value)
			arg.CalcitoninInFlush = &calcitoninInFlush
		}
		if details.Thyroglobulin.Set {
			thyroglobulin := int(details.Thyroglobulin.Value)
			arg.Thyroglobulin = &thyroglobulin
		}
	}

	// Также проверяем верхний уровень
	if req.DiagnosticMarking.Set {
		marking := domain.DiagnosticMarking(req.DiagnosticMarking.Value)
		arg.DiagnosticMarking = &marking
	}
	if req.MaterialType.Set {
		materialType := domain.MaterialType(req.MaterialType.Value)
		arg.MaterialType = &materialType
	}
	if req.Calcitonin.Set {
		calcitonin := int(req.Calcitonin.Value)
		arg.Calcitonin = &calcitonin
	}
	if req.CalcitoninInFlush.Set {
		calcitoninInFlush := int(req.CalcitoninInFlush.Value)
		arg.CalcitoninInFlush = &calcitoninInFlush
	}
	if req.Thyroglobulin.Set {
		thyroglobulin := int(req.Thyroglobulin.Value)
		arg.Thyroglobulin = &thyroglobulin
	}
	if req.IsLast.Set {
		arg.IsLast = &req.IsLast.Value
	}

	return arg
}

func (CytologyImage) UpdateArgFromCytologyUpdatePartialUpdateReq(id uuid.UUID, req *api.CytologyUpdatePartialUpdateReq) cytologySrv.UpdateCytologyImageArg {
	arg := cytologySrv.UpdateCytologyImageArg{
		Id: id,
	}

	if req.DiagnosticMarking.Set {
		marking := domain.DiagnosticMarking(req.DiagnosticMarking.Value)
		arg.DiagnosticMarking = &marking
	}
	if req.MaterialType.Set {
		materialType := domain.MaterialType(req.MaterialType.Value)
		arg.MaterialType = &materialType
	}
	if req.Calcitonin.Set {
		calcitonin := int(req.Calcitonin.Value)
		arg.Calcitonin = &calcitonin
	}
	if req.CalcitoninInFlush.Set {
		calcitoninInFlush := int(req.CalcitoninInFlush.Value)
		arg.CalcitoninInFlush = &calcitoninInFlush
	}
	if req.Thyroglobulin.Set {
		thyroglobulin := int(req.Thyroglobulin.Value)
		arg.Thyroglobulin = &thyroglobulin
	}
	if req.IsLast.Set {
		arg.IsLast = &req.IsLast.Value
	}

	return arg
}

func (CytologyImage) ToCytologyUpdateUpdateOK(img domain.CytologyImage, req *api.CytologyUpdateUpdateReq) api.CytologyUpdateUpdateOK {
	result := api.CytologyUpdateUpdateOK{
		PatientCard: api.CytologyUpdateUpdateOKPatientCard{},
		IsLast: api.OptBool{
			Value: img.IsLast,
			Set:   true,
		},
		DiagnosDate: api.OptDateTime{
			Value: img.DiagnosDate,
			Set:   true,
		},
		DiagnosticNumber: int(img.DiagnosticNumber),
	}

	if img.DiagnosticMarking != nil {
		result.DiagnosticMarking = api.OptCytologyUpdateUpdateOKDiagnosticMarking{
			Value: api.CytologyUpdateUpdateOKDiagnosticMarking(*img.DiagnosticMarking),
			Set:   true,
		}
	}

	if img.MaterialType != nil {
		result.MaterialType = api.OptCytologyUpdateUpdateOKMaterialType{
			Value: api.CytologyUpdateUpdateOKMaterialType(*img.MaterialType),
			Set:   true,
		}
	}

	if img.Calcitonin != nil {
		result.Calcitonin = api.OptInt{
			Value: int(*img.Calcitonin),
			Set:   true,
		}
	}

	if img.CalcitoninInFlush != nil {
		result.CalcitoninInFlush = api.OptInt{
			Value: int(*img.CalcitoninInFlush),
			Set:   true,
		}
	}

	if img.Thyroglobulin != nil {
		result.Thyroglobulin = api.OptInt{
			Value: int(*img.Thyroglobulin),
			Set:   true,
		}
	}

	if img.Details != nil {
		result.Details = &api.CytologyUpdateUpdateOKDetails{}
	}

	return result
}

func (CytologyImage) ToCytologyUpdatePartialUpdateOK(img domain.CytologyImage, req *api.CytologyUpdatePartialUpdateReq) api.CytologyUpdatePartialUpdateOK {
	result := api.CytologyUpdatePartialUpdateOK{
		PatientCard: api.CytologyUpdatePartialUpdateOKPatientCard{},
		IsLast: api.OptBool{
			Value: img.IsLast,
			Set:   true,
		},
		DiagnosDate: api.OptDateTime{
			Value: img.DiagnosDate,
			Set:   true,
		},
		DiagnosticNumber: int(img.DiagnosticNumber),
	}

	if img.DiagnosticMarking != nil {
		result.DiagnosticMarking = api.OptCytologyUpdatePartialUpdateOKDiagnosticMarking{
			Value: api.CytologyUpdatePartialUpdateOKDiagnosticMarking(*img.DiagnosticMarking),
			Set:   true,
		}
	}

	if img.MaterialType != nil {
		result.MaterialType = api.OptCytologyUpdatePartialUpdateOKMaterialType{
			Value: api.CytologyUpdatePartialUpdateOKMaterialType(*img.MaterialType),
			Set:   true,
		}
	}

	if img.Calcitonin != nil {
		result.Calcitonin = api.OptInt{
			Value: int(*img.Calcitonin),
			Set:   true,
		}
	}

	if img.CalcitoninInFlush != nil {
		result.CalcitoninInFlush = api.OptInt{
			Value: int(*img.CalcitoninInFlush),
			Set:   true,
		}
	}

	if img.Thyroglobulin != nil {
		result.Thyroglobulin = api.OptInt{
			Value: int(*img.Thyroglobulin),
			Set:   true,
		}
	}

	if img.Details != nil {
		result.Details = &api.CytologyUpdatePartialUpdateOKDetails{}
	}

	return result
}
