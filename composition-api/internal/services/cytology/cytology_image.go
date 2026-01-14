package cytology

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"composition-api/internal/adapters/cytology"
	domain "composition-api/internal/domain/cytology"
)

func (s *service) CreateCytologyImage(ctx context.Context, arg CreateCytologyImageArg) (uuid.UUID, error) {
	// Сначала создаем запись в БД через gRPC (без файла)
	cytologyID, err := s.adapters.Cytology.CreateCytologyImage(ctx, cytology.CreateCytologyImageIn{
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
		File:              nil, // Не передаем файл через gRPC
		ContentType:       "",
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("create cytology image in microservice: %w", err)
	}

	// Если передан файл, загружаем его в S3 напрямую (как в УЗИ)
	if arg.File != nil && (*arg.File).File != nil {
		// Генерируем ID для original_image
		originalImageID := uuid.New()

		// Формируем путь в S3: {cytology_id}/{original_image_id}/{original_image_id}
		// Используем "/" для S3, так как filepath.Join может давать разные результаты на разных ОС
		imagePath := cytologyID.String() + "/" + originalImageID.String() + "/" + originalImageID.String()

		// Загружаем файл в S3 напрямую (потоковая загрузка, без чтения в память)
		err = s.dao.NewFileRepo().LoadFile(ctx, imagePath, *arg.File)
		if err != nil {
			return uuid.Nil, fmt.Errorf("load cytology file to s3: %w", err)
		}

		// Создаем original_image через gRPC
		// Файл уже загружен в S3, но для создания записи в БД нужно передать файл через gRPC
		// Это временное решение - в будущем можно изменить протокол, чтобы передавать только путь к файлу
		// Проблема: файл все равно читается в память в адаптере CreateOriginalImage
		// Но основная проблема (таймауты при создании cytology_image) решена - файл не передается через gRPC при создании cytology_image
		_, err = s.adapters.Cytology.CreateOriginalImage(ctx, cytology.CreateOriginalImageIn{
			CytologyID:  cytologyID,
			File:        *arg.File,
			ContentType: arg.ContentType,
			DelayTime:   nil,
		})
		if err != nil {
			return uuid.Nil, fmt.Errorf("create original image: %w", err)
		}
	}

	return cytologyID, nil
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
		PrevID:            arg.PrevID,
		ParentPrevID:      arg.ParentPrevID,
	})
}

func (s *service) DeleteCytologyImage(ctx context.Context, id uuid.UUID) error {
	return s.adapters.Cytology.DeleteCytologyImage(ctx, id)
}
