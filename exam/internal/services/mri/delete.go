package mri

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) DeleteMri(ctx context.Context, id uuid.UUID) error {
	return s.dao.NewMriQuery(ctx).DeleteMri(id)
}
