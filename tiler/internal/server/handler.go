package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"tiler/internal/services"
)

type Handler struct {
	imageService services.ImageService
}

func NewHandler(imageService services.ImageService) *Handler {
	return &Handler{
		imageService: imageService,
	}
}

// GetDZI возвращает DZI XML для DeepZoom
func (h *Handler) GetDZI(w http.ResponseWriter, r *http.Request) {
	// Извлекаем путь к файлу из URL
	// Формат: /dzi/{file_path:path}
	path := strings.TrimPrefix(r.URL.Path, "/dzi/")
	if path == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	dzi, err := h.imageService.GetDZI(ctx, path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get DZI: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(dzi.XML))
}

// GetTile возвращает тайл изображения
func (h *Handler) GetTile(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из URL
	// Формат: /dzi/{file_path:path}/files/{level:int}/{col:int}_{row:int}.{format:str}
	path := r.URL.Path

	// Парсим путь
	// Пример: /dzi/path/to/image/files/5/10_20.jpeg
	parts := strings.Split(path, "/files/")
	if len(parts) != 2 {
		http.Error(w, "Invalid tile path format", http.StatusBadRequest)
		return
	}

	filePath := strings.TrimPrefix(parts[0], "/dzi/")
	tilePart := parts[1]

	// Парсим level/col_row.format
	tileParts := strings.Split(tilePart, "/")
	if len(tileParts) != 2 {
		http.Error(w, "Invalid tile path format", http.StatusBadRequest)
		return
	}

	level, err := strconv.Atoi(tileParts[0])
	if err != nil {
		http.Error(w, "Invalid level", http.StatusBadRequest)
		return
	}

	// Парсим col_row.format
	colRowFormat := strings.Split(tileParts[1], ".")
	if len(colRowFormat) != 2 {
		http.Error(w, "Invalid tile format", http.StatusBadRequest)
		return
	}

	format := colRowFormat[1]
	colRow := strings.Split(colRowFormat[0], "_")
	if len(colRow) != 2 {
		http.Error(w, "Invalid tile coordinates", http.StatusBadRequest)
		return
	}

	col, err := strconv.Atoi(colRow[0])
	if err != nil {
		http.Error(w, "Invalid column", http.StatusBadRequest)
		return
	}

	row, err := strconv.Atoi(colRow[1])
	if err != nil {
		http.Error(w, "Invalid row", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	tile, err := h.imageService.GetTile(ctx, filePath, level, col, row, format)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get tile: %v", err), http.StatusInternalServerError)
		return
	}

	// Устанавливаем Content-Type
	var contentType string
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		contentType = "image/jpeg"
	case "png":
		contentType = "image/png"
	default:
		contentType = "image/jpeg"
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(tile.Data)
}
