package cytology

import (
	"context"
	"io"

	"composition-api/internal/adapters/cytology/mappers"
	adapter_errors "composition-api/internal/adapters/errors"
	domain "composition-api/internal/domain/cytology"
	pb "composition-api/internal/generated/grpc/clients/cytology"

	"github.com/google/uuid"
)

var diagnosticMarkingMap = map[domain.DiagnosticMarking]pb.DiagnosticMarking{
	domain.DiagnosticMarkingP11: pb.DiagnosticMarking_DIAGNOSTIC_MARKING_P11,
	domain.DiagnosticMarkingL23: pb.DiagnosticMarking_DIAGNOSTIC_MARKING_L23,
}

var materialTypeMap = map[domain.MaterialType]pb.MaterialType{
	domain.MaterialTypeGS:  pb.MaterialType_MATERIAL_TYPE_GS,
	domain.MaterialTypeBP:  pb.MaterialType_MATERIAL_TYPE_BP,
	domain.MaterialTypeTP:  pb.MaterialType_MATERIAL_TYPE_TP,
	domain.MaterialTypePTP: pb.MaterialType_MATERIAL_TYPE_PTP,
	domain.MaterialTypeLNP: pb.MaterialType_MATERIAL_TYPE_LNP,
}

func (a *adapter) CreateCytologyImage(ctx context.Context, in CreateCytologyImageIn) (uuid.UUID, error) {
	req := &pb.CreateCytologyImageIn{
		ExternalId:       in.ExternalID.String(),
		DoctorId:         in.DoctorID.String(),
		PatientId:        in.PatientID.String(),
		DiagnosticNumber: int32(in.DiagnosticNumber),
	}

	// Обработка файла изображения (опционально)
	if in.File != nil && (*in.File).File != nil {
		fileData, err := io.ReadAll((*in.File).File)
		if err != nil {
			return uuid.Nil, adapter_errors.HandleGRPCError(err)
		}
		if len(fileData) > 0 {
			req.File = fileData
			if in.ContentType != "" {
				req.ContentType = &in.ContentType
			}
		}
	}

	if in.DiagnosticMarking != nil {
		dm := diagnosticMarkingMap[*in.DiagnosticMarking]
		req.DiagnosticMarking = &dm
	}

	if in.MaterialType != nil {
		mt := materialTypeMap[*in.MaterialType]
		req.MaterialType = &mt
	}

	if in.Calcitonin != nil {
		c := int32(*in.Calcitonin)
		req.Calcitonin = &c
	}

	if in.CalcitoninInFlush != nil {
		c := int32(*in.CalcitoninInFlush)
		req.CalcitoninInFlush = &c
	}

	if in.Thyroglobulin != nil {
		t := int32(*in.Thyroglobulin)
		req.Thyroglobulin = &t
	}

	if in.Details != nil {
		req.Details = in.Details
	}

	if in.PrevID != nil {
		prev := in.PrevID.String()
		req.PrevId = &prev
	}

	if in.ParentPrevID != nil {
		parent := in.ParentPrevID.String()
		req.ParentPrevId = &parent
	}

	res, err := a.client.CreateCytologyImage(ctx, req)
	if err != nil {
		return uuid.Nil, adapter_errors.HandleGRPCError(err)
	}

	return uuid.MustParse(res.Id), nil
}

func (a *adapter) GetCytologyImageById(ctx context.Context, id uuid.UUID) (domain.CytologyImage, *domain.OriginalImage, error) {
	res, err := a.client.GetCytologyImageById(ctx, &pb.GetCytologyImageByIdIn{Id: id.String()})
	if err != nil {
		return domain.CytologyImage{}, nil, adapter_errors.HandleGRPCError(err)
	}

	img := mappers.CytologyImage{}.Domain(res.CytologyImage)

	var originalImage *domain.OriginalImage
	if res.OriginalImage != nil {
		oi := mappers.OriginalImage{}.Domain(res.OriginalImage)
		originalImage = &oi
	}

	return img, originalImage, nil
}

func (a *adapter) GetCytologyImagesByExternalId(ctx context.Context, id uuid.UUID) ([]domain.CytologyImage, error) {
	res, err := a.client.GetCytologyImagesByExternalId(ctx, &pb.GetCytologyImagesByExternalIdIn{ExternalId: id.String()})
	if err != nil {
		return nil, adapter_errors.HandleGRPCError(err)
	}

	return mappers.CytologyImage{}.SliceDomain(res.CytologyImages), nil
}

func (a *adapter) GetCytologyImagesByDoctorIdAndPatientId(ctx context.Context, doctorID, patientID uuid.UUID) ([]domain.CytologyImage, error) {
	res, err := a.client.GetCytologyImagesByDoctorIdAndPatientId(ctx, &pb.GetCytologyImagesByDoctorIdAndPatientIdIn{
		DoctorId:  doctorID.String(),
		PatientId: patientID.String(),
	})
	if err != nil {
		return nil, adapter_errors.HandleGRPCError(err)
	}

	return mappers.CytologyImage{}.SliceDomain(res.CytologyImages), nil
}

func (a *adapter) UpdateCytologyImage(ctx context.Context, in UpdateCytologyImageIn) (domain.CytologyImage, error) {
	req := &pb.UpdateCytologyImageIn{
		Id: in.Id.String(),
	}

	if in.DiagnosticMarking != nil {
		dm := diagnosticMarkingMap[*in.DiagnosticMarking]
		req.DiagnosticMarking = &dm
	}

	if in.MaterialType != nil {
		mt := materialTypeMap[*in.MaterialType]
		req.MaterialType = &mt
	}

	if in.Calcitonin != nil {
		c := int32(*in.Calcitonin)
		req.Calcitonin = &c
	}

	if in.CalcitoninInFlush != nil {
		c := int32(*in.CalcitoninInFlush)
		req.CalcitoninInFlush = &c
	}

	if in.Thyroglobulin != nil {
		t := int32(*in.Thyroglobulin)
		req.Thyroglobulin = &t
	}

	if in.Details != nil {
		req.Details = in.Details
	}

	if in.IsLast != nil {
		req.IsLast = in.IsLast
	}

	if in.PrevID != nil {
		prevIDStr := in.PrevID.String()
		req.PrevId = &prevIDStr
	}

	if in.ParentPrevID != nil {
		parentPrevIDStr := in.ParentPrevID.String()
		req.ParentPrevId = &parentPrevIDStr
	}

	res, err := a.client.UpdateCytologyImage(ctx, req)
	if err != nil {
		return domain.CytologyImage{}, adapter_errors.HandleGRPCError(err)
	}

	return mappers.CytologyImage{}.Domain(res.CytologyImage), nil
}

func (a *adapter) DeleteCytologyImage(ctx context.Context, id uuid.UUID) error {
	_, err := a.client.DeleteCytologyImage(ctx, &pb.DeleteCytologyImageIn{Id: id.String()})
	return adapter_errors.HandleGRPCError(err)
}

func (a *adapter) CopyCytologyImage(ctx context.Context, id uuid.UUID) (domain.CytologyImage, error) {
	res, err := a.client.CopyCytologyImage(ctx, &pb.CopyCytologyImageIn{Id: id.String()})
	if err != nil {
		return domain.CytologyImage{}, adapter_errors.HandleGRPCError(err)
	}

	return mappers.CytologyImage{}.Domain(res.CytologyImage), nil
}

func (a *adapter) GetCytologyImageHistory(ctx context.Context, id uuid.UUID) ([]domain.CytologyImage, error) {
	res, err := a.client.GetCytologyImageHistory(ctx, &pb.GetCytologyImageHistoryIn{Id: id.String()})
	if err != nil {
		return nil, adapter_errors.HandleGRPCError(err)
	}

	return mappers.CytologyImage{}.SliceDomain(res.CytologyImages), nil
}
