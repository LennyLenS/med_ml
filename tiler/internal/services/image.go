package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

type cachedImage struct {
	path      string
	image     *vips.ImageRef
	lastUsed  time.Time
	accessMux sync.RWMutex
}

type imageService struct {
	s3Client     S3Client
	tileSize     int
	overlap      int
	cache        map[string]*cachedImage
	cacheMux     sync.RWMutex
	cacheDir     string
	maxCacheSize int
	cacheTTL     time.Duration
}

type S3Client interface {
	GetObject(ctx context.Context, bucketName, objectName string) ([]byte, error)
}

func NewImageService(s3Client S3Client, tileSize, overlap int) ImageService {
	// Инициализируем libvips (требуется только один раз)
	vips.Startup(nil)

	// Создаем временную директорию для кэша файлов
	cacheDir := filepath.Join(os.TempDir(), "tiler_cache")
	os.MkdirAll(cacheDir, 0755)

	return &imageService{
		s3Client:     s3Client,
		tileSize:     tileSize,
		overlap:      overlap,
		cache:        make(map[string]*cachedImage),
		cacheDir:     cacheDir,
		maxCacheSize: 10,        // Максимум 10 открытых файлов в кэше
		cacheTTL:     time.Hour, // Время жизни кэша
	}
}

func (s *imageService) GetImageInfo(ctx context.Context, imagePath string) (*domain.ImageInfo, error) {
	// Используем кэшированное изображение для получения метаданных
	// libvips читает только заголовок файла, не весь файл
	img, err := s.getOrLoadImage(ctx, imagePath)
	if err != nil {
		return nil, err
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

// getOrLoadImage получает изображение из кэша или загружает его
func (s *imageService) getOrLoadImage(ctx context.Context, imagePath string) (*vips.ImageRef, error) {
	// Проверяем кэш
	s.cacheMux.RLock()
	cached, exists := s.cache[imagePath]
	s.cacheMux.RUnlock()

	if exists {
		cached.accessMux.Lock()
		cached.lastUsed = time.Now()
		// Создаем копию изображения для использования
		imgCopy, err := cached.image.Copy()
		cached.accessMux.Unlock()
		if err == nil {
			return imgCopy, nil
		}
		// Если копирование не удалось, удаляем из кэша и загружаем заново
		s.cacheMux.Lock()
		delete(s.cache, imagePath)
		s.cacheMux.Unlock()
	}

	// Файл не в кэше, загружаем его
	localPath := filepath.Join(s.cacheDir, strings.ReplaceAll(imagePath, "/", "_"))

	// Проверяем, существует ли файл локально
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		// Скачиваем файл из S3
		imgData, err := s.s3Client.GetObject(ctx, "", imagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to get image from S3 (path: %s): %w", imagePath, err)
		}

		// Сохраняем во временный файл
		if err := os.WriteFile(localPath, imgData, 0644); err != nil {
			return nil, fmt.Errorf("failed to save image to cache: %w", err)
		}
		// Освобождаем память
		imgData = nil
	}

	// Открываем файл через libvips с random access
	// libvips автоматически будет читать только нужные тайлы из tiled TIFF
	vipsImg, err := vips.NewImageFromFile(localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image with libvips: %w", err)
	}

	// Добавляем в кэш
	s.cacheMux.Lock()
	// Очищаем старые записи, если кэш переполнен
	if len(s.cache) >= s.maxCacheSize {
		s.cleanupCache()
	}
	s.cache[imagePath] = &cachedImage{
		path:     localPath,
		image:    vipsImg,
		lastUsed: time.Now(),
	}
	s.cacheMux.Unlock()

	// Возвращаем копию для использования
	return vipsImg.Copy()
}

// cleanupCache удаляет старые записи из кэша
func (s *imageService) cleanupCache() {
	now := time.Now()
	for path, cached := range s.cache {
		if now.Sub(cached.lastUsed) > s.cacheTTL {
			cached.accessMux.Lock()
			cached.image.Close()
			cached.accessMux.Unlock()
			delete(s.cache, path)
		}
	}
}

func (s *imageService) GetTile(ctx context.Context, imagePath string, level, col, row int, format string) (*domain.Tile, error) {
	// Сначала получаем информацию об изображении для валидации
	info, err := s.GetImageInfo(ctx, imagePath)
	if err != nil {
		return nil, err
	}

	// Проверяем, что уровень валиден
	if level < 0 || level >= info.Levels {
		return nil, fmt.Errorf("invalid level: %d (max: %d)", level, info.Levels-1)
	}

	// Вычисляем размеры изображения на данном уровне
	var levelWidth, levelHeight int
	if level == 0 {
		levelWidth = info.Width
		levelHeight = info.Height
	} else {
		scale := math.Pow(2, float64(level))
		levelWidth = int(float64(info.Width) / scale)
		levelHeight = int(float64(info.Height) / scale)
	}

	// Проверяем, что размеры изображения валидны
	if levelWidth <= 0 || levelHeight <= 0 {
		return nil, fmt.Errorf("invalid image dimensions at level %d: %dx%d", level, levelWidth, levelHeight)
	}

	// Вычисляем максимальное количество тайлов на данном уровне
	// Даже если изображение меньше тайла, должен быть хотя бы один тайл
	maxCol := int(math.Max(1, math.Ceil(float64(levelWidth)/float64(s.tileSize))))
	maxRow := int(math.Max(1, math.Ceil(float64(levelHeight)/float64(s.tileSize))))

	// Проверяем границы координат тайла
	if col < 0 || row < 0 || col >= maxCol || row >= maxRow {
		return nil, fmt.Errorf("tile coordinates out of bounds: level=%d, col=%d (max=%d), row=%d (max=%d), image_size=%dx%d",
			level, col, maxCol-1, row, maxRow-1, levelWidth, levelHeight)
	}

	// Получаем изображение из кэша или загружаем его
	// libvips будет читать только нужные тайлы из файла благодаря random access
	vipsImg, err := s.getOrLoadImage(ctx, imagePath)
	if err != nil {
		return nil, err
	}
	defer vipsImg.Close()

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

	// Дополнительная проверка границ после масштабирования (на случай округления)
	if x >= vipsImg.Width() || y >= vipsImg.Height() {
		return nil, errors.New("tile coordinates out of bounds")
	}

	// Обрезаем тайл с учетом overlap
	cropX := int(math.Max(0, float64(x-s.overlap)))
	cropY := int(math.Max(0, float64(y-s.overlap)))
	cropWidth := int(math.Min(float64(vipsImg.Width()-cropX), float64(s.tileSize+2*s.overlap)))
	cropHeight := int(math.Min(float64(vipsImg.Height()-cropY), float64(s.tileSize+2*s.overlap)))

	// Извлекаем тайл
	// ExtractArea модифицирует исходное изображение, поэтому нужно создать копию
	var tileImg *vips.ImageRef
	if cropX > 0 || cropY > 0 || cropWidth < vipsImg.Width() || cropHeight < vipsImg.Height() {
		// Создаем копию изображения для извлечения области
		tileImg, err = vipsImg.Copy()
		if err != nil {
			return nil, fmt.Errorf("failed to copy image: %w", err)
		}
		defer tileImg.Close()

		if err := tileImg.ExtractArea(cropX, cropY, cropWidth, cropHeight); err != nil {
			return nil, fmt.Errorf("failed to extract tile: %w", err)
		}
	} else {
		// Если не нужно обрезать, используем исходное изображение
		tileImg = vipsImg
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
