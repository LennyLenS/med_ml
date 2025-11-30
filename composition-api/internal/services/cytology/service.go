package cytology

import (
	"context"

	"github.com/google/uuid"

	"composition-api/internal/adapters"
	domain "composition-api/internal/domain/cytology"
)

type Service interface {
	CreateCytologyImage(ctx context.Context, arg CreateCytologyImageArg) (uuid.UUID, error)
	GetCytologyImageById(ctx context.Context, id uuid.UUID) (domain.CytologyImage, *domain.OriginalImage, error)
	GetCytologyImagesByExternalId(ctx context.Context, externalID uuid.UUID) ([]domain.CytologyImage, error)
	GetCytologyImagesByPatientCardId(ctx context.Context, patientCardID uuid.UUID) ([]domain.CytologyImage, error)
	UpdateCytologyImage(ctx context.Context, arg UpdateCytologyImageArg) (domain.CytologyImage, error)
	DeleteCytologyImage(ctx context.Context, id uuid.UUID) error

	CreateOriginalImage(ctx context.Context, arg CreateOriginalImageArg) (uuid.UUID, error)
	GetOriginalImageById(ctx context.Context, id uuid.UUID) (domain.OriginalImage, error)
	GetOriginalImagesByCytologyId(ctx context.Context, id uuid.UUID) ([]domain.OriginalImage, error)
	UpdateOriginalImage(ctx context.Context, arg UpdateOriginalImageArg) (domain.OriginalImage, error)

	CreateSegmentationGroup(ctx context.Context, arg CreateSegmentationGroupArg) (uuid.UUID, error)
	GetSegmentationGroupsByCytologyId(ctx context.Context, id uuid.UUID, segType *domain.SegType, groupType *domain.GroupType, isAI *bool) ([]domain.SegmentationGroup, error)
	UpdateSegmentationGroup(ctx context.Context, arg UpdateSegmentationGroupArg) (domain.SegmentationGroup, error)
	DeleteSegmentationGroup(ctx context.Context, id uuid.UUID) error

	CreateSegmentation(ctx context.Context, arg CreateSegmentationArg) (uuid.UUID, error)
	GetSegmentationById(ctx context.Context, id uuid.UUID) (domain.Segmentation, error)
	GetSegmentsByGroupId(ctx context.Context, id uuid.UUID) ([]domain.Segmentation, error)
	UpdateSegmentation(ctx context.Context, arg UpdateSegmentationArg) (domain.Segmentation, error)
	DeleteSegmentation(ctx context.Context, id uuid.UUID) error
}

type service struct {
	adapters *adapters.Adapters
}

func New(adapters *adapters.Adapters) Service {
	return &service{
		adapters: adapters,
	}
}
