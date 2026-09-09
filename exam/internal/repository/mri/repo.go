package mri

import (
	"exam/internal/repository/mri/entity"

	daolib "github.com/WantBeASleep/med_ml_lib/dao"
	"github.com/google/uuid"
)

const (
	table = "mri"

	columnID          = "id"
	columnProjection  = "projection"
	columnChecked     = "checked"
	columnExternalID  = "external_id"
	columnAuthor      = "author"
	columnDeviceID    = "device_id"
	columnStatus      = "status"
	columnDescription = "description"
	columnCreateAt    = "create_at"
)

type Repository interface {
	CheckExist(id uuid.UUID) (bool, error)

	InsertMri(mri entity.Mri) error

	GetMriByID(id uuid.UUID) (entity.Mri, error)
	GetMrisByExternalID(externalID uuid.UUID) ([]entity.Mri, error)
	GetMrisByAuthor(author uuid.UUID) ([]entity.Mri, error)

	UpdateMri(mri entity.Mri) error
	UpdateMriStatus(id uuid.UUID, status string) error

	DeleteMri(id uuid.UUID) error
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
