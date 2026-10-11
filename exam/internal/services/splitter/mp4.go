package splitter

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"exam/internal/domain"
)

func (s *service) DicomToMp4(ctx context.Context, file domain.File) (domain.File, error) {
	if err := ctx.Err(); err != nil {
		return domain.File{}, err
	}
	frames, err := decodeDicomFrames(file, false)
	if err != nil {
		return domain.File{}, err
	}
	if len(frames) < 54 {
		return domain.File{}, fmt.Errorf("CT model requires at least 54 DICOM frames, got %d", len(frames))
	}
	dir, err := os.MkdirTemp("", "exam-dicom-*")
	if err != nil {
		return domain.File{}, err
	}
	defer os.RemoveAll(dir)
	input, err := os.Create(filepath.Join(dir, "frames.png"))
	if err != nil {
		return domain.File{}, err
	}
	for _, frame := range frames {
		if err := ctx.Err(); err != nil {
			input.Close()
			return domain.File{}, err
		}
		if err := png.Encode(input, frame); err != nil {
			input.Close()
			return domain.File{}, fmt.Errorf("encode video frame: %w", err)
		}
	}
	if err := input.Close(); err != nil {
		return domain.File{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	output := filepath.Join(dir, "converted.mp4")
	cmd := exec.CommandContext(ctx, "ffmpeg", "-nostdin", "-v", "error", "-f", "image2pipe", "-framerate", "15", "-i", input.Name(), "-an", "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2", "-c:v", "mpeg4", "-q:v", "2", "-pix_fmt", "yuv420p", "-movflags", "+faststart", output)
	if err := cmd.Run(); err != nil {
		return domain.File{}, fmt.Errorf("convert dicom to mp4 with ffmpeg: %w", err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		return domain.File{}, err
	}
	if len(data) == 0 {
		return domain.File{}, fmt.Errorf("ffmpeg produced an empty mp4")
	}
	return domain.File{Format: Mp4, Size: int64(len(data)), Buf: bytes.NewReader(data)}, nil
}
