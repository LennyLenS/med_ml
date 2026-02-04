package splitter

import (
	"fmt"
	"image/jpeg"

	"uzi/internal/domain"
)

type jpegSplitter struct{}

func (j jpegSplitter) splitToPng(f domain.File) ([]domain.File, error) {
	// JPEG файлы обычно содержат один кадр
	// Декодируем JPEG и конвертируем в PNG
	img, err := jpeg.Decode(f.Buf)
	if err != nil {
		return nil, fmt.Errorf("decode jpeg: %w", err)
	}

	// Конвертируем в PNG
	pngFile, err := convertToPng(img)
	if err != nil {
		return nil, fmt.Errorf("convert to png: %w", err)
	}

	return []domain.File{pngFile}, nil
}
