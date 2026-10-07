package cytology_analysis

import (
	"context"
	"time"

	"github.com/google/uuid"

	"cytology/internal/repository"
	analysisentity "cytology/internal/repository/cytology_analysis/entity"
)

type Analysis struct {
	ID              uuid.UUID
	CytologyID      uuid.UUID
	OriginalImageID uuid.UUID
}

type Service interface {
	Create(ctx context.Context, cytologyID, originalImageID uuid.UUID) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (Analysis, error)
	Complete(ctx context.Context, id uuid.UUID, result []byte) error
}

type service struct {
	dao repository.DAO
}

func New(dao repository.DAO) Service {
	return &service{dao: dao}
}

func (s *service) Create(ctx context.Context, cytologyID, originalImageID uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()
	err := s.dao.NewCytologyAnalysisQuery(ctx).Insert(analysisentity.CytologyAnalysis{
		ID:              id,
		CytologyID:      cytologyID,
		OriginalImageID: originalImageID,
		CreatedAt:       time.Now(),
	})
	return id, err
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (Analysis, error) {
	analysis, err := s.dao.NewCytologyAnalysisQuery(ctx).GetByID(id)
	if err != nil {
		return Analysis{}, err
	}
	return Analysis{ID: analysis.ID, CytologyID: analysis.CytologyID, OriginalImageID: analysis.OriginalImageID}, nil
}

func (s *service) Complete(ctx context.Context, id uuid.UUID, result []byte) error {
	return s.dao.NewCytologyAnalysisQuery(ctx).Complete(id, result)
}
