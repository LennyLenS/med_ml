package kt

import (
	"github.com/google/uuid"
	ht "github.com/ogen-go/ogen/http"
)

type CreateKtArg struct {
	File        ht.MultipartFile
	Author      uuid.UUID
	DeviceID    int
	Description *string
}

type UpdateKtArg struct {
	Id                 uuid.UUID
	Checked            *bool
	ClassProbabilities map[string]float64
}
