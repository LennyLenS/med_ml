package image

import (
	"exam/internal/repository/image/entity"

	daolib "github.com/WantBeASleep/med_ml_lib/dao"
	"github.com/google/uuid"
)

const (
	table = "image"

	columnId    = "id"
	columnMriId = "mri_id"
	columnPage  = "page"
)

type Repository interface {
	InsertImages(images ...entity.Image) error

	GetImagesByMriID(mriID uuid.UUID) ([]entity.Image, error)
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
