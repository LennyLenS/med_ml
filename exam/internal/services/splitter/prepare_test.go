package splitter

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"exam/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type storedFile struct {
	format string
	data   []byte
}

type memoryFiles struct {
	files    map[string]storedFile
	failPath string
}

func (r *memoryFiles) GetFileViaTemp(context.Context, string) (domain.File, func() error, error) {
	return domain.File{}, nil, errors.New("unexpected read")
}

func (r *memoryFiles) LoadFile(_ context.Context, path string, file domain.File) error {
	if path == r.failPath {
		return errors.New("storage failure")
	}
	data, err := io.ReadAll(file.Buf)
	r.files[path] = storedFile{file.Format, data}
	return err
}

func TestPrepareDicomPreservesOriginalAndRewindsTiff(t *testing.T) {
	id := uuid.New()
	data := dicomFixture(1, 2, 2, 1, 8, 8, 0, 0, "MONOCHROME2", []byte{0, 255, 128, 64})
	repo := &memoryFiles{files: make(map[string]storedFile)}
	file, err := PrepareDicom(context.Background(), repo, id, domain.File{Format: Dicom, Size: int64(len(data)), Buf: bytes.NewReader(data)}, Tiff)
	require.NoError(t, err)
	require.Equal(t, storedFile{Dicom, data}, repo.files[id.String()+"/original.dcm"])
	encoded, err := io.ReadAll(file.Buf)
	require.NoError(t, err)
	require.Equal(t, storedFile{Tiff, encoded}, repo.files[id.String()+"/"+id.String()])
	pages, err := New().SplitToPng(domain.File{Format: Tiff, Buf: bytes.NewReader(encoded)})
	require.NoError(t, err)
	require.Len(t, pages, 2)
}

func TestPrepareDicomRejectsBadInputWithoutWriting(t *testing.T) {
	repo := &memoryFiles{files: make(map[string]storedFile)}
	_, err := PrepareDicom(context.Background(), repo, uuid.New(), domain.File{Format: Dicom, Buf: bytes.NewReader([]byte("broken"))}, Tiff)
	require.Error(t, err)
	require.Empty(t, repo.files)
}

func TestPrepareDicomStopsOnStorageFailure(t *testing.T) {
	id := uuid.New()
	data := dicomFixture(1, 2, 1, 1, 8, 8, 0, 0, "MONOCHROME2", []byte{0, 255})
	repo := &memoryFiles{files: make(map[string]storedFile), failPath: id.String() + "/original.dcm"}
	_, err := PrepareDicom(context.Background(), repo, id, domain.File{Format: Dicom, Buf: bytes.NewReader(data)}, Tiff)
	require.ErrorContains(t, err, "preserve original")
	require.Empty(t, repo.files)
}

func TestPrepareDicomLeavesExistingFormatsUntouched(t *testing.T) {
	for _, format := range []string{Tiff, Mp4} {
		input := domain.File{Format: format, Buf: bytes.NewReader([]byte("unchanged"))}
		repo := &memoryFiles{files: make(map[string]storedFile)}
		output, err := PrepareDicom(context.Background(), repo, uuid.New(), input, format)
		require.NoError(t, err)
		require.Equal(t, input, output)
		require.Empty(t, repo.files)
	}
}
