package splitter

import (
	"uzi/internal/domain"
)

type pngSplitter struct{}

func (pngSplitter) splitToPng(f domain.File) ([]domain.File, error) {
	// PNG файлы уже в нужном формате, просто возвращаем как есть
	// Убеждаемся, что формат установлен правильно
	result := f
	result.Format = Png
	return []domain.File{result}, nil
}
