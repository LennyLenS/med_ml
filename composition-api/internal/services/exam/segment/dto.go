package segment

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateSegmentArg struct {
	ImageID   uuid.UUID
	NodeID    uuid.UUID
	Contor    json.RawMessage
	Knosp_012 float64
	Knosp_3   float64
	Knosp_4   float64
}

type UpdateSegmentArg struct {
	Id        uuid.UUID
	Contor    json.RawMessage
	Knosp_012 *float64
	Knosp_3   *float64
	Knosp_4   *float64
}
