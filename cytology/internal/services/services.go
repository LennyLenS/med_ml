package services

import (
	"context"

	dbus "cytology/internal/dbus/producers"
	"cytology/internal/repository"
	"cytology/internal/services/cytology_analysis"
	"cytology/internal/services/cytology_image"
	"cytology/internal/services/original_image"
	"cytology/internal/services/segmentation"
	"cytology/internal/services/segmentation_group"
)

type Services struct {
	CytologyImage     cytology_image.Service
	CytologyAnalysis  cytology_analysis.Service
	OriginalImage     original_image.Service
	SegmentationGroup segmentation_group.Service
	Segmentation      segmentation.Service
	Producer          dbus.Producer
	dao               repository.DAO
}

func New(
	dao repository.DAO,
	producer dbus.Producer,
) *Services {
	cytologyImage := cytology_image.New(dao)
	cytologyAnalysis := cytology_analysis.New(dao)
	originalImage := original_image.New(dao)
	segmentationGroup := segmentation_group.New(dao)
	segmentation := segmentation.New(dao)

	return &Services{
		CytologyImage:     cytologyImage,
		CytologyAnalysis:  cytologyAnalysis,
		OriginalImage:     originalImage,
		SegmentationGroup: segmentationGroup,
		Segmentation:      segmentation,
		Producer:          producer,
		dao:               dao,
	}
}

func (s *Services) BeginTx(ctx context.Context) (context.Context, error) {
	return s.dao.BeginTx(ctx)
}

func (s *Services) CommitTx(ctx context.Context) error {
	return s.dao.CommitTx(ctx)
}

func (s *Services) RollbackTx(ctx context.Context) error {
	return s.dao.RollbackTx(ctx)
}
