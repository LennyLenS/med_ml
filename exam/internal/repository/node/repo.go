package node

import (
	"exam/internal/repository/node/entity"

	daolib "github.com/WantBeASleep/med_ml_lib/dao"
	"github.com/google/uuid"
)

const (
	table = "node"

	columnID          = "id"
	columnAI          = "ai"
	columnMriID       = "mri_id"
	columnValidation  = "validation"
	columnKnosp012    = "knosp_012"
	columnKnosp3      = "knosp_3"
	columnKnosp4      = "knosp_4"
	columnDescription = "description"
)

type Repository interface {
	InsertNodes(nodes ...entity.Node) error

	GetNodeByID(id uuid.UUID) (entity.Node, error)
	GetNodesByImageID(id uuid.UUID) ([]entity.Node, error)
	GetNodesByMriID(id uuid.UUID) ([]entity.Node, error)

	UpdateNode(node entity.Node) error

	DeleteNodeByID(id uuid.UUID) error
}

type repo struct {
	*daolib.BaseQuery
}

func NewRepo() *repo {
	return &repo{}
}

func (q *repo) SetBaseQuery(baseQuery *daolib.BaseQuery) {
	q.BaseQuery = baseQuery
}
