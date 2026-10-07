package segment

import (
	"encoding/json"

	"github.com/google/uuid"

	"exam/internal/domain"
)

type CreateSegmentArg struct {
	ImageID  uuid.UUID
	NodeID   uuid.UUID
	Contor   json.RawMessage
	Knosp012 float64
	Knosp3   float64
	Knosp4   float64
}

// TODO: починить баг при запросе со всеми полями nil
type UpdateSegmentArg struct {
	Id       uuid.UUID
	Contor   *json.RawMessage
	Knosp012 *float64
	Knosp3   *float64
	Knosp4   *float64
}

func (u UpdateSegmentArg) UpdateDomain(d *domain.Segment) {
	if u.Contor != nil {
		d.Contor = *u.Contor
	}
	if u.Knosp012 != nil {
		d.Knosp012 = *u.Knosp012
	}
	if u.Knosp3 != nil {
		d.Knosp3 = *u.Knosp3
	}
	if u.Knosp4 != nil {
		d.Knosp4 = *u.Knosp4
	}
}
