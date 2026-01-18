package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"math"
	"strings"

	"tiler/internal/domain"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/tiff" // Регистрирует TIFF декодер
)

type ImageService interface {
	GetImageInfo(ctx context.Context, imagePath string) (*domain.ImageInfo, error)
	GetTile(ctx context.Context, imagePath string, level, col, row int, format string) (*domain.Tile, error)
	GetDZI(ctx context.Context, imagePath string) (*domain.DZI, error)
}

type imageService struct {
	s3Client S3Client
	tileSize int
	overlap  int
}

type S3Client interface {
	GetObject(ctx context.Context, bucketName, objectName string) ([]byte, error)
}

func NewImageService(s3Client S3Client, tileSize, overlap int) ImageService {
	return &imageService{
		s3Client: s3Client,
		tileSize: tileSize,
		overlap:  overlap,
	}
}

func (s *imageService) GetImageInfo(ctx context.Context, imagePath string) (*domain.ImageInfo, error) {
	// Загружаем изображение из S3
	// Для больших файлов (3GB) это может быть проблематично, но для получения размеров
	// используем DecodeConfig, который читает только метаданные, а не декодирует весь файл
	imgData, err := s.s3Client.GetObject(ctx, "", imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from S3 (path: %s): %w", imagePath, err)
	}

	// Используем DecodeConfig вместо Decode - он читает только метаданные (размеры)
	// и не декодирует весь файл в память, что критично для больших SVS файлов
	cfg, _, err := image.DecodeConfig(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image config (compression value 7 may not be supported): %w", err)
	}

	width := cfg.Width
	height := cfg.Height

	// Вычисляем количество уровней (levels)
	// Уровень 0 - оригинальное изображение, каждый следующий уровень в 2 раза меньше
	maxDim := math.Max(float64(width), float64(height))
	levels := int(math.Ceil(math.Log2(maxDim/float64(s.tileSize)))) + 1

	return &domain.ImageInfo{
		Width:    width,
		Height:   height,
		Levels:   levels,
		TileSize: s.tileSize,
		Overlap:  s.overlap,
	}, nil
}

func (s *imageService) GetTile(ctx context.Context, imagePath string, level, col, row int, format string) (*domain.Tile, error) {
	// Загружаем изображение из S3
	// imagePath должен быть путем внутри bucket (например: "cytology_id/image_id/image_id")
	imgData, err := s.s3Client.GetObject(ctx, "", imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from S3 (path: %s): %w", imagePath, err)
	}

	// Декодируем изображение
	// Используем image.Decode, который автоматически выберет правильный декодер
	// golang.org/x/image/tiff должен поддерживать compression value 7 в новых версиях
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Масштабируем изображение до нужного уровня
	scaledImg := s.scaleToLevel(img, level)

	// Вычисляем размеры тайла с учетом overlap
	tileSizeWithOverlap := s.tileSize + 2*s.overlap

	// Вычисляем координаты тайла
	x := col * s.tileSize
	y := row * s.tileSize

	// Проверяем границы
	bounds := scaledImg.Bounds()
	if x >= bounds.Dx() || y >= bounds.Dy() {
		return nil, errors.New("tile coordinates out of bounds")
	}

	// Вычисляем размеры для обрезки
	width := s.tileSize
	height := s.tileSize
	if x+width > bounds.Dx() {
		width = bounds.Dx() - x
	}
	if y+height > bounds.Dy() {
		height = bounds.Dy() - y
	}

	// Обрезаем тайл с учетом overlap
	cropX := math.Max(0, float64(x-s.overlap))
	cropY := math.Max(0, float64(y-s.overlap))
	cropWidth := math.Min(float64(bounds.Dx()-int(cropX)), float64(tileSizeWithOverlap))
	cropHeight := math.Min(float64(bounds.Dy()-int(cropY)), float64(tileSizeWithOverlap))

	tileImg := imaging.Crop(scaledImg, image.Rect(
		int(cropX),
		int(cropY),
		int(cropX)+int(cropWidth),
		int(cropY)+int(cropHeight),
	))

	// Кодируем тайл в нужный формат
	var tileData []byte
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		tileData, err = encodeJPEG(tileImg)
	case "png":
		tileData, err = encodePNG(tileImg)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode tile: %w", err)
	}

	return &domain.Tile{
		Level:  level,
		Col:    col,
		Row:    row,
		Data:   tileData,
		Format: format,
	}, nil
}

func (s *imageService) scaleToLevel(img image.Image, level int) image.Image {
	if level == 0 {
		return img
	}

	// Вычисляем масштаб для уровня
	scale := math.Pow(2, float64(level))
	bounds := img.Bounds()
	newWidth := int(float64(bounds.Dx()) / scale)
	newHeight := int(float64(bounds.Dy()) / scale)

	return imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
}

func (s *imageService) GetDZI(ctx context.Context, imagePath string) (*domain.DZI, error) {
	info, err := s.GetImageInfo(ctx, imagePath)
	if err != nil {
		return nil, err
	}

	// Генерируем DZI XML
	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Image xmlns="http://schemas.microsoft.com/deepzoom/2008"
       TileSize="%d"
       Overlap="%d"
       Format="jpeg"
       ServerFormat="Default">
  <Size Width="%d" Height="%d"/>
</Image>`, s.tileSize, s.overlap, info.Width, info.Height)

	return &domain.DZI{
		XML:       xml,
		ImageInfo: *info,
	}, nil
}

func encodeJPEG(img image.Image) ([]byte, error) {
	// Используем imaging для кодирования
	var buf bytes.Buffer
	err := imaging.Encode(&buf, img, imaging.JPEG)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := imaging.Encode(&buf, img, imaging.PNG)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
