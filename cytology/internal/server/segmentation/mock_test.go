package segmentation_test

import (
	"context"

	"cytology/internal/domain"
	"cytology/internal/server/segmentation"
	segmentationservice "cytology/internal/services/segmentation"
	"cytology/internal/services"
)

type mockSegmentationService struct {
	createID  int
	createErr error
	seg       domain.Segmentation
	segErr    error
	segs      []domain.Segmentation
	segsErr   error
	updateSeg domain.Segmentation
	updateErr error
	deleteErr error
}

func (m *mockSegmentationService) CreateSegmentation(context.Context, segmentationservice.CreateSegmentationArg) (int, error) {
	return m.createID, m.createErr
}

func (m *mockSegmentationService) GetSegmentationByID(context.Context, int) (domain.Segmentation, error) {
	return m.seg, m.segErr
}

func (m *mockSegmentationService) GetSegmentsByGroupID(context.Context, int) ([]domain.Segmentation, error) {
	return m.segs, m.segsErr
}

func (m *mockSegmentationService) UpdateSegmentation(context.Context, segmentationservice.UpdateSegmentationArg) (domain.Segmentation, error) {
	return m.updateSeg, m.updateErr
}

func (m *mockSegmentationService) DeleteSegmentation(context.Context, int) error {
	return m.deleteErr
}

func newHandler(svc segmentationservice.Service) segmentation.SegmentationHandler {
	return segmentation.New(&services.Services{Segmentation: svc})
}
