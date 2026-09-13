package node_segment

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateNodesWithSegmentsArgNode struct {
	Knosp012    float64
	Knosp3      float64
	Knosp4      float64
	Description *string
}

type CreateNodesWithSegmentsArgSegment struct {
	ImageID  uuid.UUID
	Contor   json.RawMessage
	Knosp012 float64
	Knosp3   float64
	Knosp4   float64
}

type CreateNodesWithSegmentsArg struct {
	Node     CreateNodesWithSegmentsArgNode
	Segments []CreateNodesWithSegmentsArgSegment
}

type CreateNodesWithSegmentsID struct {
	NodeID     uuid.UUID
	SegmentsID []uuid.UUID
}
