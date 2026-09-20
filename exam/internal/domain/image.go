package domain

import (
	"github.com/google/uuid"
)

type Image struct {
	Id    uuid.UUID
	MriID uuid.UUID
	Page  int
}
