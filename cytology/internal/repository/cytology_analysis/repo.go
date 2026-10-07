package cytology_analysis

import (
	daolib "github.com/WantBeASleep/med_ml_lib/dao"
	"github.com/google/uuid"

	"cytology/internal/repository/cytology_analysis/entity"
)

const table = "cytology_analysis"

type Repository interface {
	Insert(analysis entity.CytologyAnalysis) error
	GetByID(id uuid.UUID) (entity.CytologyAnalysis, error)
	Complete(id uuid.UUID, result []byte) error
}

type repo struct {
	*daolib.BaseQuery
}

func NewR() *repo {
	return &repo{}
}

func (q *repo) SetBaseQuery(baseQuery *daolib.BaseQuery) {
	q.BaseQuery = baseQuery
}
