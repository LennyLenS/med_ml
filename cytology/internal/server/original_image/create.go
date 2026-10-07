package original_image

import (
	"context"
	"fmt"

	"cytology/internal/domain"
	cytologyanalysisrequestedpb "cytology/internal/generated/dbus/produce/cytologyanalysisrequested"
	pb "cytology/internal/generated/grpc/service"
	"cytology/internal/services/original_image"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handler) CreateOriginalImage(ctx context.Context, in *pb.CreateOriginalImageIn) (*pb.CreateOriginalImageOut, error) {
	cytologyID, err := uuid.Parse(in.CytologyId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cytology_id is not a valid uuid: %s", err.Error())
	}

	// Если передан путь к файлу, используем его (файл уже загружен в S3)
	// Иначе проверяем, что передан файл
	imagePath := in.GetImagePath()
	var imagePathPtr *string
	if imagePath != "" {
		imagePathPtr = &imagePath
	}

	if imagePathPtr == nil || *imagePathPtr == "" {
		if len(in.File) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "file or image_path is required")
		}
		if in.ContentType == "" {
			return nil, status.Errorf(codes.InvalidArgument, "content_type is required")
		}
	}

	var delayTimePtr *float64
	if in.DelayTime != nil {
		delayTime := in.GetDelayTime()
		delayTimePtr = &delayTime
	}

	arg := original_image.CreateOriginalImageArg{
		CytologyID:  cytologyID,
		File:        in.File,
		ContentType: in.ContentType,
		DelayTime:   delayTimePtr,
		ImagePath:   imagePathPtr,
	}

	id, err := h.services.OriginalImage.CreateOriginalImage(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	if h.services.Producer != nil {
		originalImage, err := h.services.OriginalImage.GetOriginalImageByID(ctx, id)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "get original image: %s", err.Error())
		}
		cytologyImage, err := h.services.CytologyImage.GetCytologyImageByID(ctx, cytologyID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "get cytology image: %s", err.Error())
		}

		analysisID, err := h.services.CytologyAnalysis.Create(ctx, cytologyID, id)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "create cytology analysis: %s", err.Error())
		}

		materialTypes := map[domain.MaterialType]cytologyanalysisrequestedpb.MaterialType{
			domain.MaterialTypeGS:  cytologyanalysisrequestedpb.MaterialType_MATERIAL_TYPE_GS,
			domain.MaterialTypeBP:  cytologyanalysisrequestedpb.MaterialType_MATERIAL_TYPE_BP,
			domain.MaterialTypeTP:  cytologyanalysisrequestedpb.MaterialType_MATERIAL_TYPE_TP,
			domain.MaterialTypePTP: cytologyanalysisrequestedpb.MaterialType_MATERIAL_TYPE_PTP,
			domain.MaterialTypeLNP: cytologyanalysisrequestedpb.MaterialType_MATERIAL_TYPE_LNP,
		}
		materialType := cytologyanalysisrequestedpb.MaterialType_MATERIAL_TYPE_UNSPECIFIED
		if cytologyImage.MaterialType != nil {
			materialType = materialTypes[*cytologyImage.MaterialType]
		}

		msg := &cytologyanalysisrequestedpb.AnalysisRequested{
			AnalysisId:       analysisID.String(),
			InputArtifactUri: fmt.Sprintf("s3://cytology/%s", originalImage.ImagePath),
			MaterialType:     materialType,
		}
		if cytologyImage.MaterialType != nil && *cytologyImage.MaterialType == domain.MaterialTypeLNP {
			if cytologyImage.CalcitoninInFlush != nil {
				value := int32(*cytologyImage.CalcitoninInFlush)
				msg.CalcitoninInFlush = &value
			}
			if cytologyImage.Thyroglobulin != nil {
				value := int32(*cytologyImage.Thyroglobulin)
				msg.ThyroglobulinInFlush = &value
			}
		}
		if err := h.services.Producer.SendCytologyAnalysisRequested(ctx, msg); err != nil {
			return nil, status.Errorf(codes.Internal, "send cytology analysis requested: %s", err.Error())
		}
	}

	return &pb.CreateOriginalImageOut{Id: id.String()}, nil
}
