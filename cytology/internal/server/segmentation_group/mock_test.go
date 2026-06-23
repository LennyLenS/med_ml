package segmentation_group_test

import (
	"context"

	"cytology/internal/domain"
	"cytology/internal/server/segmentation_group"
	segmentationgroupservice "cytology/internal/services/segmentation_group"
	"cytology/internal/services"

	"github.com/google/uuid"
)

type mockSegmentationGroupService struct {
	createID  int
	createErr error
	groups    []domain.SegmentationGroup
	groupsErr error
	group     domain.SegmentationGroup
	groupErr  error
	updateErr error
	deleteErr error
}

func (m *mockSegmentationGroupService) CreateSegmentationGroup(context.Context, segmentationgroupservice.CreateSegmentationGroupArg) (int, error) {
	return m.createID, m.createErr
}

func (m *mockSegmentationGroupService) GetSegmentationGroupByID(context.Context, int) (domain.SegmentationGroup, error) {
	return m.group, m.groupErr
}

func (m *mockSegmentationGroupService) GetSegmentationGroupsByCytologyID(context.Context, uuid.UUID) ([]domain.SegmentationGroup, error) {
	return m.groups, m.groupsErr
}

func (m *mockSegmentationGroupService) UpdateSegmentationGroup(context.Context, segmentationgroupservice.UpdateSegmentationGroupArg) (domain.SegmentationGroup, error) {
	if m.updateErr != nil {
		return domain.SegmentationGroup{}, m.updateErr
	}
	return m.group, nil
}

func (m *mockSegmentationGroupService) DeleteSegmentationGroup(context.Context, int) error {
	return m.deleteErr
}

func newHandler(svc segmentationgroupservice.Service) segmentation_group.SegmentationGroupHandler {
	return segmentation_group.New(&services.Services{SegmentationGroup: svc})
}
