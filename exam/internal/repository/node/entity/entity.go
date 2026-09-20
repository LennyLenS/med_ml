package entity

import (
	"database/sql"

	"github.com/WantBeASleep/med_ml_lib/gtc"
	"github.com/google/uuid"

	"exam/internal/domain"
)

type Node struct {
	Id          uuid.UUID      `db:"id"`
	Ai          bool           `db:"ai"`
	MriID       uuid.UUID      `db:"mri_id"`
	Validation  sql.NullString `db:"validation"`
	Knosp012    float64        `db:"knosp_012"`
	Knosp3      float64        `db:"knosp_3"`
	Knosp4      float64        `db:"knosp_4"`
	Description sql.NullString `db:"description"`
}

func (Node) FromDomain(d domain.Node) Node {
	return Node{
		Id:          d.Id,
		Ai:          d.Ai,
		MriID:       d.MriID,
		Validation:  gtc.String.PointerToSql((*string)(d.Validation)),
		Knosp012:    d.Knosp012,
		Knosp3:      d.Knosp3,
		Knosp4:      d.Knosp4,
		Description: gtc.String.PointerToSql(d.Description),
	}
}

func (Node) SliceFromDomain(slice []domain.Node) []Node {
	res := make([]Node, 0, len(slice))
	for _, v := range slice {
		res = append(res, Node{}.FromDomain(v))
	}
	return res
}

func (d Node) ToDomain() domain.Node {
	node := domain.Node{
		Id:          d.Id,
		Ai:          d.Ai,
		MriID:       d.MriID,
		Knosp012:    d.Knosp012,
		Knosp3:      d.Knosp3,
		Knosp4:      d.Knosp4,
		Description: gtc.String.SqlToPointer(d.Description),
	}

	if d.Validation.Valid {
		validation, _ := domain.NodeValidation.Parse("", d.Validation.String)
		node.Validation = &validation
	}

	return node
}

func (Node) SliceToDomain(slice []Node) []domain.Node {
	res := make([]domain.Node, 0, len(slice))
	for _, v := range slice {
		res = append(res, v.ToDomain())
	}
	return res
}
