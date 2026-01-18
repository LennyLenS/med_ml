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

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/disintegration/imaging"
	_ "golang.org/x/image/tiff" // Регистрирует TIFF декодер для fallback
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
	// Инициализируем libvips (требуется только один раз)
	vips.Startup(nil)

	return &imageService{
		s3Client: s3Client,
		tileSize: tileSize,
		overlap:  overlap,
	}
}

func (s *imageService) GetImageInfo(ctx context.Context, imagePath string) (*domain.ImageInfo, error) {
	// Загружаем изображение из S3
	imgData, err := s.s3Client.GetObject(ctx, "", imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from S3 (path: %s): %w", imagePath, err)
	}

	// Используем govips для получения метаданных
	// govips поддерживает SVS файлы и compression value 7
	img, err := vips.NewImageFromBuffer(imgData)
	if err != nil {
		// Fallback на стандартный декодер для совместимости
		cfg, _, decodeErr := image.DecodeConfig(bytes.NewReader(imgData))
		if decodeErr != nil {
			return nil, fmt.Errorf("failed to decode image (govips error: %v, standard decoder error: %w)", err, decodeErr)
		}
		width := cfg.Width
		height := cfg.Height
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
	defer img.Close()

	width := img.Width()
	height := img.Height()

	// Вычисляем количество уровней (levels)
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
	// Загружаем файл из S3
	imgData, err := s.s3Client.GetObject(ctx, "", imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from S3 (path: %s): %w", imagePath, err)
	}

	// Используем govips для декодирования - поддерживает SVS и compression value 7
	vipsImg, err := vips.NewImageFromBuffer(imgData)
	if err != nil {
		// Fallback на стандартный декодер
		img, _, decodeErr := image.Decode(bytes.NewReader(imgData))
		if decodeErr != nil {
			return nil, fmt.Errorf("failed to decode image (govips error: %v, standard decoder error: %w)", err, decodeErr)
		}
		return s.getTileFromStandardImage(img, level, col, row, format)
	}
	defer vipsImg.Close()

	// Освобождаем память как можно скорее
	imgData = nil

	// Масштабируем изображение до нужного уровня
	if level > 0 {
		scale := math.Pow(2, float64(level))
		scaleFactor := 1.0 / scale
		if err := vipsImg.Resize(scaleFactor, vips.KernelLanczos3); err != nil {
			return nil, fmt.Errorf("failed to resize image: %w", err)
		}
	}

	// Вычисляем координаты тайла
	x := col * s.tileSize
	y := row * s.tileSize

	// Проверяем границы
	if x >= vipsImg.Width() || y >= vipsImg.Height() {
		return nil, errors.New("tile coordinates out of bounds")
	}

	// Обрезаем тайл с учетом overlap
	cropX := int(math.Max(0, float64(x-s.overlap)))
	cropY := int(math.Max(0, float64(y-s.overlap)))
	cropWidth := int(math.Min(float64(vipsImg.Width()-cropX), float64(s.tileSize+2*s.overlap)))
	cropHeight := int(math.Min(float64(vipsImg.Height()-cropY), float64(s.tileSize+2*s.overlap)))

	// Извлекаем тайл
	tileImg := vipsImg
	if cropX > 0 || cropY > 0 || cropWidth < vipsImg.Width() || cropHeight < vipsImg.Height() {
		if err := tileImg.ExtractArea(cropX, cropY, cropWidth, cropHeight); err != nil {
			return nil, fmt.Errorf("failed to extract tile: %w", err)
		}
	}

	// Кодируем тайл в нужный формат
	var tileData []byte
	var encodeErr error
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		ep := vips.NewJpegExportParams()
		tileData, _, encodeErr = tileImg.ExportJpeg(ep)
	case "png":
		ep := vips.NewPngExportParams()
		tileData, _, encodeErr = tileImg.ExportPng(ep)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if encodeErr != nil {
		return nil, fmt.Errorf("failed to encode tile: %w", encodeErr)
	}

	return &domain.Tile{
		Level:  level,
		Col:    col,
		Row:    row,
		Data:   tileData,
		Format: format,
	}, nil
}

// getTileFromStandardImage - fallback метод для стандартного image.Decode
func (s *imageService) getTileFromStandardImage(img image.Image, level, col, row int, format string) (*domain.Tile, error) {
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
	var encodeErr error
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		tileData, encodeErr = encodeJPEG(tileImg)
	case "png":
		tileData, encodeErr = encodePNG(tileImg)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if encodeErr != nil {
		return nil, fmt.Errorf("failed to encode tile: %w", encodeErr)
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
