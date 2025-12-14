package download

import (
	"context"
	"io"
	"path/filepath"

	"github.com/google/uuid"
)

func (s *service) GetCytologyImage(ctx context.Context, cytologyID uuid.UUID, originalImageID uuid.UUID) (io.ReadCloser, error) {
	return s.repo.NewFileRepo().GetFile(
		ctx,
		filepath.Join(
			cytologyID.String(),
			originalImageID.String(),
			originalImageID.String(),
		),
	)
}
