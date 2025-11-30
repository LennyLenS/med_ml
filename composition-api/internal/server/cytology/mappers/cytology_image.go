package mappers

import (
	"crypto/sha256"

	"github.com/google/uuid"

	api "composition-api/internal/generated/http/api"
	domain "composition-api/internal/domain/cytology"
	cytologySrv "composition-api/internal/services/cytology"
)

// GetPatientCardID generates a deterministic UUID from doctorID and patientID
func GetPatientCardID(doctorID, patientID uuid.UUID) uuid.UUID {
	hash := sha256.New()
	hash.Write(doctorID[:])
	hash.Write(patientID[:])
	sum := hash.Sum(nil)

	var uuidBytes [16]byte
	copy(uuidBytes[:], sum[:16])

	uuidBytes[6] = (uuidBytes[6] & 0x0f) | 0x40 // Version 4
	uuidBytes[8] = (uuidBytes[8] & 0x3f) | 0x80 // Variant 10

	return uuid.UUID(uuidBytes)
}

type CytologyImage struct{}

func (CytologyImage) Domain(img domain.CytologyImage, doctorID, patientID uuid.UUID) api.CytologyImage {
	result := api.CytologyImage{
		ID:               img.Id,
		ExternalID:       img.ExternalID,
		DoctorID:         doctorID,
		PatientID:        patientID,
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

func (CytologyImage) SliceDomain(imgs []domain.CytologyImage, doctorID, patientID uuid.UUID) []api.CytologyImage {
	result := make([]api.CytologyImage, 0, len(imgs))
	for _, img := range imgs {
		result = append(result, CytologyImage{}.Domain(img, doctorID, patientID))
	}
	return result
}

func (CytologyImage) CreateArg(req *api.CytologyPostReq) cytologySrv.CreateCytologyImageArg {
	patientCardID := GetPatientCardID(req.DoctorID, req.PatientID)

	arg := cytologySrv.CreateCytologyImageArg{
		ExternalID:       req.ExternalID,
		PatientCardID:    patientCardID,
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

func (CytologyImage) UpdateArg(id uuid.UUID, req *api.CytologyIDPatchReq) cytologySrv.UpdateCytologyImageArg {
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
