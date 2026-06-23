package original_image_test

import (
	"context"

	"cytology/internal/domain"
	"cytology/internal/server/original_image"
	originalimageservice "cytology/internal/services/original_image"
	"cytology/internal/services"

	"github.com/google/uuid"
)

type mockOriginalImageService struct {
	createID    uuid.UUID
	createErr   error
	image       domain.OriginalImage
	imageErr    error
	images      []domain.OriginalImage
	imagesErr   error
	updateImage domain.OriginalImage
	updateErr   error
}

func (m *mockOriginalImageService) CreateOriginalImage(context.Context, originalimageservice.CreateOriginalImageArg) (uuid.UUID, error) {
	return m.createID, m.createErr
}

func (m *mockOriginalImageService) GetOriginalImageByID(context.Context, uuid.UUID) (domain.OriginalImage, error) {
	return m.image, m.imageErr
}

func (m *mockOriginalImageService) GetOriginalImagesByCytologyID(context.Context, uuid.UUID) ([]domain.OriginalImage, error) {
	return m.images, m.imagesErr
}

func (m *mockOriginalImageService) UpdateOriginalImage(context.Context, originalimageservice.UpdateOriginalImageArg) (domain.OriginalImage, error) {
	return m.updateImage, m.updateErr
}

func newHandler(svc originalimageservice.Service) original_image.OriginalImageHandler {
	return original_image.New(&services.Services{OriginalImage: svc})
}
