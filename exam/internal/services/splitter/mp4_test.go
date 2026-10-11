package splitter

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"exam/internal/domain"

	"github.com/stretchr/testify/require"
)

func TestDicomToMp4Frames(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg is required")
	}
	for _, count := range []int{54, 60} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			pixels := make([]byte, count*16*16)
			for i := range pixels {
				pixels[i] = uint8(40 + i/(16*16)*3)
			}
			data := dicomFixture(16, 16, count, 1, 8, 8, 0, 0, "MONOCHROME2", pixels)
			converted, err := New().DicomToMp4(context.Background(), domain.File{Format: Dicom, Size: int64(len(data)), Buf: bytes.NewReader(data)})
			require.NoError(t, err)
			require.Equal(t, Mp4, converted.Format)
			video, err := io.ReadAll(converted.Buf)
			require.NoError(t, err)
			require.Equal(t, int64(len(video)), converted.Size)
			path := filepath.Join(t.TempDir(), "video.mp4")
			require.NoError(t, os.WriteFile(path, video, 0600))
			decoded, err := exec.Command("ffmpeg", "-v", "error", "-i", path, "-f", "rawvideo", "-pix_fmt", "gray", "pipe:1").Output()
			require.NoError(t, err)
			require.Len(t, decoded, len(pixels))
			for i, value := range decoded {
				require.InDelta(t, pixels[i], value, 4)
			}
		})
	}
}

func TestDicomToMp4RejectsInsufficientCTFrames(t *testing.T) {
	data := dicomFixture(1, 2, 1, 1, 8, 8, 0, 0, "MONOCHROME2", []byte{0, 255})
	_, err := New().DicomToMp4(context.Background(), domain.File{Format: Dicom, Buf: bytes.NewReader(data)})
	require.ErrorContains(t, err, "at least 54")
}

func TestDicomToMp4Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	data := dicomFixture(1, 2, 1, 1, 8, 8, 0, 0, "MONOCHROME2", []byte{0, 255})
	_, err := New().DicomToMp4(ctx, domain.File{Format: Dicom, Buf: bytes.NewReader(data)})
	require.ErrorIs(t, err, context.Canceled)
}
