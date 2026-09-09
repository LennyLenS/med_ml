package mri

import (
	domain "composition-api/internal/domain/exam"

	"github.com/google/uuid"
	ht "github.com/ogen-go/ogen/http"
)

type CreateMriArg struct {
	File        ht.MultipartFile
	Projection  domain.MriProjection
	ExternalID  uuid.UUID
	Author      uuid.UUID
	DeviceID    int
	Description *string
}

type UpdateMriArg struct {
	Id         uuid.UUID
	Projection *domain.MriProjection
	Checked    *bool
}
