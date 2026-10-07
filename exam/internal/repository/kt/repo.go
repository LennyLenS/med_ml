package kt

import (
	"exam/internal/repository/kt/entity"

	daolib "github.com/WantBeASleep/med_ml_lib/dao"
	"github.com/google/uuid"
)

const (
	table = "kt"

	columnID               = "id"
	columnChecked          = "checked"
	columnAuthor           = "author"
	columnDeviceID         = "device_id"
	columnStatus           = "status"
	columnDescription      = "description"
	columnCreateAt         = "create_at"
	columnPredictedClasses = "predicted_classes"
)

type Repository interface {
	CheckExist(id uuid.UUID) (bool, error)

	InsertKt(kt entity.KT) error

	GetKtByID(id uuid.UUID) (entity.KT, error)

	GetKtsByAuthor(id uuid.UUID) ([]entity.KT, error)

	UpdateKt(mri entity.KT) error

	DeleteKt(id uuid.UUID) error

	UpdateKtPrediction(id uuid.UUID, predictedClasses []byte) error
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
