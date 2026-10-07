package domain

import "github.com/google/uuid"

type Node struct {
	Id          uuid.UUID
	Ai          bool
	MriID       uuid.UUID
	Validation  *NodeValidation
	Knosp012    float64
	Knosp3      float64
	Knosp4      float64
	Description *string
}
