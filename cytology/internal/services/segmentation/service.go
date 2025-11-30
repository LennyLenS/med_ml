package segmentation

import (
	"context"
	"time"

	"cytology/internal/domain"
	segmentationEntity "cytology/internal/repository/segmentation/entity"
	"cytology/internal/repository"

	"github.com/google/uuid"
)

type Service interface {
	CreateSegmentation(ctx context.Context, arg CreateSegmentationArg) (uuid.UUID, error)
	GetSegmentationByID(ctx context.Context, id uuid.UUID) (domain.Segmentation, error)
	GetSegmentsByGroupID(ctx context.Context, groupID uuid.UUID) ([]domain.Segmentation, error)
	UpdateSegmentation(ctx context.Context, arg UpdateSegmentationArg) (domain.Segmentation, error)
	DeleteSegmentation(ctx context.Context, id uuid.UUID) error
}

type service struct {
	dao repository.DAO
}

func New(dao repository.DAO) Service {
	return &service{
		dao: dao,
	}
}

type CreateSegmentationArg struct {
	SegmentationGroupID uuid.UUID
	Points              []domain.SegmentationPoint
}

type UpdateSegmentationArg struct {
	Id     uuid.UUID
	Points []domain.SegmentationPoint
}

func (s *service) CreateSegmentation(ctx context.Context, arg CreateSegmentationArg) (uuid.UUID, error) {
	seg := domain.Segmentation{
		Id:                uuid.New(),
		SegmentationGroupID: arg.SegmentationGroupID,
		Points:            arg.Points,
		CreateAt:          time.Now(),
	}

	// Генерируем ID для точек
	for i := range seg.Points {
		seg.Points[i].Id = uuid.New()
		seg.Points[i].SegmentationID = seg.Id
		seg.Points[i].CreateAt = time.Now()
	}

	entitySeg := segmentationEntity.Segmentation{}.FromDomain(seg)
	if err := s.dao.NewSegmentationQuery(ctx).InsertSegmentation(entitySeg); err != nil {
		return uuid.Nil, err
	}

	return seg.Id, nil
}

func (s *service) GetSegmentationByID(ctx context.Context, id uuid.UUID) (domain.Segmentation, error) {
	seg, err := s.dao.NewSegmentationQuery(ctx).GetSegmentationByID(id)
	if err != nil {
		return domain.Segmentation{}, err
	}
	return seg.ToDomain(), nil
}

func (s *service) GetSegmentsByGroupID(ctx context.Context, groupID uuid.UUID) ([]domain.Segmentation, error) {
	segs, err := s.dao.NewSegmentationQuery(ctx).GetSegmentsByGroupID(groupID)
	if err != nil {
		return nil, err
	}
	return segmentationEntity.Segmentation{}.SliceToDomain(segs), nil
}

func (s *service) UpdateSegmentation(ctx context.Context, arg UpdateSegmentationArg) (domain.Segmentation, error) {
	seg, err := s.dao.NewSegmentationQuery(ctx).GetSegmentationByID(arg.Id)
	if err != nil {
		return domain.Segmentation{}, err
	}

	domainSeg := seg.ToDomain()
	domainSeg.Points = arg.Points

	// Обновляем ID и время создания для точек
	for i := range domainSeg.Points {
		domainSeg.Points[i].Id = uuid.New()
		domainSeg.Points[i].SegmentationID = domainSeg.Id
		domainSeg.Points[i].CreateAt = time.Now()
	}

	entitySeg := segmentationEntity.Segmentation{}.FromDomain(domainSeg)
	if err := s.dao.NewSegmentationQuery(ctx).UpdateSegmentation(entitySeg); err != nil {
		return domain.Segmentation{}, err
	}

	return domainSeg, nil
}

func (s *service) DeleteSegmentation(ctx context.Context, id uuid.UUID) error {
	return s.dao.NewSegmentationQuery(ctx).DeleteSegmentation(id)
}
