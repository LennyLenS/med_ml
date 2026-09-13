package node

import (
	domain "composition-api/internal/domain/exam"

	"github.com/google/uuid"
)

type UpdateNodeArg struct {
	Id         uuid.UUID
	Validation *domain.NodeValidation
	Knosp_012  *float64
	Knosp_3    *float64
	Knosp_4    *float64
}
