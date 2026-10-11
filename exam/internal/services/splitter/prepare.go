package splitter

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"exam/internal/domain"
	"exam/internal/repository"

	"github.com/google/uuid"
)

func PrepareDicom(ctx context.Context, repo repository.FileRepo, id uuid.UUID, file domain.File, target string) (domain.File, error) {
	if file.Format != Dicom {
		return file, nil
	}
	var converted domain.File
	var err error
	s := New()
	switch target {
	case Tiff:
		converted, err = s.DicomToTiff(file)
	case Mp4:
		converted, err = s.DicomToMp4(ctx, file)
	default:
		return domain.File{}, ErrUnsupportedFormat
	}
	if err != nil {
		return domain.File{}, err
	}
	original, ok := file.Buf.(io.Seeker)
	if !ok {
		return domain.File{}, fmt.Errorf("dicom input is not seekable")
	}
	if _, err := original.Seek(0, io.SeekStart); err != nil {
		return domain.File{}, err
	}
	if err := repo.LoadFile(ctx, filepath.Join(id.String(), "original.dcm"), file); err != nil {
		return domain.File{}, fmt.Errorf("preserve original dicom: %w", err)
	}
	if err := repo.LoadFile(ctx, filepath.Join(id.String(), id.String()), converted); err != nil {
		return domain.File{}, fmt.Errorf("save converted file: %w", err)
	}
	if _, err := converted.Buf.(io.Seeker).Seek(0, io.SeekStart); err != nil {
		return domain.File{}, err
	}
	return converted, nil
}
