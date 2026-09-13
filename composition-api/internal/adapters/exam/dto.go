package exam

import (
	"encoding/json"

	"github.com/google/uuid"

	domain "composition-api/internal/domain/exam"
)

type CreateMriIn struct {
	Projection  domain.MriProjection
	ExternalID  uuid.UUID
	Author      uuid.UUID
	DeviceID    int
	Description *string
}

type CreateKtIn struct {
	Author      uuid.UUID
	DeviceID    int
	Description *string
}

type UpdateKtIn struct {
	Id                 uuid.UUID
	Checked            *bool
	ClassProbabilities map[string]float64
}

type UpdateMriIn struct {
	Id         uuid.UUID
	Projection *domain.MriProjection
	Checked    *bool
}

type UpdateNodeIn struct {
	Id         uuid.UUID
	Validation *domain.NodeValidation
	Knosp_012  *float64
	Knosp_3    *float64
	Knosp_4    *float64
}

type CreateSegmentIn struct {
	ImageID   uuid.UUID
	NodeID    uuid.UUID
	Contor    json.RawMessage
	Knosp_012 float64
	Knosp_3   float64
	Knosp_4   float64
}

type UpdateSegmentIn struct {
	Id        uuid.UUID
	Contor    json.RawMessage
	Knosp_012 *float64
	Knosp_3   *float64
	Knosp_4   *float64
}

type CreateNodeWithSegmentsIn_Node struct {
	Knosp_012   float64
	Knosp_3     float64
	Knosp_4     float64
	Description *string
}

type CreateNodeWithSegmentsIn_Segment struct {
	ImageID   uuid.UUID
	Contor    json.RawMessage
	Knosp_012 float64
	Knosp_3   float64
	Knosp_4   float64
}

type CreateNodeWithSegmentsIn struct {
	MriID    uuid.UUID
	Node     CreateNodeWithSegmentsIn_Node
	Segments []CreateNodeWithSegmentsIn_Segment
}


