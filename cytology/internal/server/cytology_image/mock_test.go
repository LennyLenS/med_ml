package cytology_image_test

import (
	"context"

	"cytology/internal/domain"
	"cytology/internal/server/cytology_image"
	cytologyimageservice "cytology/internal/services/cytology_image"
	originalimageservice "cytology/internal/services/original_image"
	"cytology/internal/services"

	"github.com/google/uuid"
)

type mockCytologyImageService struct {
	createID    uuid.UUID
	createErr   error
	image       domain.CytologyImage
	imageErr    error
	images      []domain.CytologyImage
	imagesErr   error
	ids         []uuid.UUID
	idsErr      error
	updateImage domain.CytologyImage
	updateErr   error
	deleteErr   error
	copyImage   domain.CytologyImage
	copyErr     error
	history     []domain.CytologyImage
	historyErr  error
}

func (m *mockCytologyImageService) CreateCytologyImage(context.Context, cytologyimageservice.CreateCytologyImageArg) (uuid.UUID, error) {
	return m.createID, m.createErr
}

func (m *mockCytologyImageService) GetCytologyImageByID(context.Context, uuid.UUID) (domain.CytologyImage, error) {
	return m.image, m.imageErr
}

func (m *mockCytologyImageService) GetCytologyImagesByExternalID(context.Context, uuid.UUID) ([]domain.CytologyImage, error) {
	return m.images, m.imagesErr
}

func (m *mockCytologyImageService) GetCytologyImagesByDoctorIdAndPatientId(context.Context, uuid.UUID, uuid.UUID) ([]domain.CytologyImage, error) {
	return m.images, m.imagesErr
}

func (m *mockCytologyImageService) GetCytologyImagesByPatientId(context.Context, uuid.UUID) ([]domain.CytologyImage, error) {
	return m.images, m.imagesErr
}

func (m *mockCytologyImageService) GetCytologyImageIdsByDoctorIdAndPatientId(context.Context, uuid.UUID, uuid.UUID) ([]uuid.UUID, error) {
	return m.ids, m.idsErr
}

func (m *mockCytologyImageService) UpdateCytologyImage(context.Context, cytologyimageservice.UpdateCytologyImageArg) (domain.CytologyImage, error) {
	return m.updateImage, m.updateErr
}

func (m *mockCytologyImageService) DeleteCytologyImage(context.Context, uuid.UUID) error {
	return m.deleteErr
}

func (m *mockCytologyImageService) CopyCytologyImage(context.Context, uuid.UUID) (domain.CytologyImage, error) {
	return m.copyImage, m.copyErr
}

func (m *mockCytologyImageService) GetCytologyImageHistory(context.Context, uuid.UUID) ([]domain.CytologyImage, error) {
	return m.history, m.historyErr
}

type mockOriginalImageService struct {
	images []domain.OriginalImage
	err    error
}

func (m *mockOriginalImageService) CreateOriginalImage(context.Context, originalimageservice.CreateOriginalImageArg) (uuid.UUID, error) {
	panic("not implemented")
}

func (m *mockOriginalImageService) GetOriginalImageByID(context.Context, uuid.UUID) (domain.OriginalImage, error) {
	panic("not implemented")
}

func (m *mockOriginalImageService) GetOriginalImagesByCytologyID(context.Context, uuid.UUID) ([]domain.OriginalImage, error) {
	return m.images, m.err
}

func (m *mockOriginalImageService) UpdateOriginalImage(context.Context, originalimageservice.UpdateOriginalImageArg) (domain.OriginalImage, error) {
	panic("not implemented")
}

func newHandler(svc cytologyimageservice.Service, originalSvc ...originalimageservice.Service) cytology_image.CytologyImageHandler {
	s := &services.Services{CytologyImage: svc}
	if len(originalSvc) > 0 {
		s.OriginalImage = originalSvc[0]
	}
	return cytology_image.New(s)
}
