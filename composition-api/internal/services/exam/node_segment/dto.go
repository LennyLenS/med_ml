package node_segment

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateNodeWithSegmentArg_Node struct {
	Knosp_012   float64
	Knosp_3     float64
	Knosp_4     float64
	Description *string
}

type CreateNodeWithSegmentArg_Segment struct {
	ImageID   uuid.UUID
	Contor    json.RawMessage
	Knosp_012 float64
	Knosp_3   float64
	Knosp_4   float64
}

type CreateNodeWithSegmentArg struct {
	MriID    uuid.UUID
	Node     CreateNodeWithSegmentArg_Node
	Segments []CreateNodeWithSegmentArg_Segment
}
