package domain

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Segment struct {
	Id       uuid.UUID
	ImageID  uuid.UUID
	NodeID   uuid.UUID
	Contor   json.RawMessage
	Ai       bool
	Knosp012 float64
	Knosp3   float64
	Knosp4   float64
}
