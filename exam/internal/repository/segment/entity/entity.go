package entity

import (
	"encoding/json"

	"exam/internal/domain"

	"github.com/google/uuid"
)

type Segment struct {
	Id       uuid.UUID       `db:"id"`
	ImageID  uuid.UUID       `db:"image_id"`
	NodeID   uuid.UUID       `db:"node_id"`
	Contor   json.RawMessage `db:"contor"`
	Ai       bool            `db:"ai"`
	Knosp012 float64         `db:"knosp_012"`
	Knosp3   float64         `db:"knosp_3"`
	Knosp4   float64         `db:"knosp_4"`
}

func (Segment) FromDomain(d domain.Segment) Segment {
	return Segment{
		Id:       d.Id,
		ImageID:  d.ImageID,
		NodeID:   d.NodeID,
		Contor:   d.Contor,
		Ai:       d.Ai,
		Knosp012: d.Knosp012,
		Knosp3:   d.Knosp3,
		Knosp4:   d.Knosp4,
	}
}

func (Segment) SliceFromDomain(slice []domain.Segment) []Segment {
	res := make([]Segment, 0, len(slice))
	for _, v := range slice {
		res = append(res, Segment{}.FromDomain(v))
	}
	return res
}

func (d Segment) ToDomain() domain.Segment {
	return domain.Segment{
		Id:       d.Id,
		ImageID:  d.ImageID,
		NodeID:   d.NodeID,
		Contor:   d.Contor,
		Ai:       d.Ai,
		Knosp012: d.Knosp012,
		Knosp3:   d.Knosp3,
		Knosp4:   d.Knosp4,
	}
}

func (Segment) SliceToDomain(slice []Segment) []domain.Segment {
	res := make([]domain.Segment, 0, len(slice))
	for _, v := range slice {
		res = append(res, v.ToDomain())
	}
	return res
}
