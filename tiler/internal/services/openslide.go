//go:build cgo
// +build cgo

package services

/*
#cgo pkg-config: openslide
#include <openslide.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unsafe"

	"tiler/internal/domain"

	"github.com/disintegration/imaging"
)

type openslideImage struct {
	osr       *C.openslide_t
	path      string
	lastUsed  time.Time
	accessMux sync.RWMutex
}

type openslideService struct {
	s3Client     S3Client
	tileSize     int
	overlap      int
	cache        map[string]*openslideImage
	cacheMux     sync.RWMutex
	infoCache    map[string]*domain.ImageInfo
	infoCacheMux sync.RWMutex
	cacheDir     string
	tileCacheDir string
	maxCacheSize int
	cacheTTL     time.Duration
}

func NewOpenSlideService(s3Client S3Client, tileSize, overlap int) ImageService {
	// Проверяем, что OpenSlide доступен
	version := C.openslide_get_version()
	if version == nil {
		slog.Error("OpenSlide library not found - check if openslide-dev is installed")
		return nil
	}
	slog.Info("OpenSlide version", "version", C.GoString(version))

	cacheDir := filepath.Join(os.TempDir(), "tiler_cache")
	os.MkdirAll(cacheDir, 0755)

	tileCacheDir := filepath.Join(os.TempDir(), "tiler_tile_cache")
	os.MkdirAll(tileCacheDir, 0755)

	return &openslideService{
		s3Client:     s3Client,
		tileSize:     tileSize,
		overlap:      overlap,
		cache:        make(map[string]*openslideImage),
		infoCache:    make(map[string]*domain.ImageInfo),
		cacheDir:     cacheDir,
		tileCacheDir: tileCacheDir,
		maxCacheSize: 10,
		cacheTTL:     time.Hour,
	}
}

func (s *openslideService) GetImageInfo(ctx context.Context, imagePath string) (*domain.ImageInfo, error) {
	// Проверяем кэш
	s.infoCacheMux.RLock()
	cachedInfo, exists := s.infoCache[imagePath]
	s.infoCacheMux.RUnlock()

	if exists && cachedInfo != nil {
		return &domain.ImageInfo{
			Width:    cachedInfo.Width,
			Height:   cachedInfo.Height,
			Levels:   cachedInfo.Levels,
			TileSize: cachedInfo.TileSize,
			Overlap:  cachedInfo.Overlap,
		}, nil
	}

	// Загружаем изображение для получения информации
	osImg, err := s.getOrLoadImage(ctx, imagePath)
	if err != nil {
		return nil, err
	}
	defer s.closeImage(osImg)

	// Получаем размеры уровня 0 (полное разрешение)
	var width, height C.int64_t
	C.openslide_get_level0_dimensions(osImg.osr, &width, &height)

	// Получаем количество уровней
	levelCount := int(C.openslide_get_level_count(osImg.osr))

	// Вычисляем количество уровней для DZI
	// Формула: levels = ceil(log2(max(width, height) / tileSize)) + 1
	maxDim := math.Max(float64(width), float64(height))
	levels := int(math.Ceil(math.Log2(maxDim/float64(s.tileSize)))) + 1

	// Используем максимум из встроенных уровней OpenSlide и вычисленных для DZI
	if levelCount > levels {
		levels = levelCount
	}

	info := &domain.ImageInfo{
		Width:    int(width),
		Height:   int(height),
		Levels:   levels,
		TileSize: s.tileSize,
		Overlap:  s.overlap,
	}

	// Сохраняем в кэш
	s.infoCacheMux.Lock()
	s.infoCache[imagePath] = info
	s.infoCacheMux.Unlock()

	return info, nil
}

func (s *openslideService) GetTile(ctx context.Context, imagePath string, level, col, row int, format string) (*domain.Tile, error) {
	info, err := s.GetImageInfo(ctx, imagePath)
	if err != nil {
		return nil, err
	}

	// Проверяем валидность уровня
	if level < 0 || level >= info.Levels {
		return nil, fmt.Errorf("invalid level: %d (max: %d)", level, info.Levels-1)
	}

	// Вычисляем размеры изображения на данном уровне
	// В OpenSeadragon: maxLevel = полное разрешение, 0 = минимальный масштаб
	// Формула: scaleFactor = 2^(level - maxLevel)
	maxLevel := info.Levels - 1
	scaleFactor := math.Pow(2, float64(level-maxLevel))
	levelWidth := int(float64(info.Width) * scaleFactor)
	levelHeight := int(float64(info.Height) * scaleFactor)

	// Проверяем валидность размеров
	if levelWidth <= 0 || levelHeight <= 0 {
		return nil, fmt.Errorf("invalid level dimensions: level=%d, width=%d, height=%d (maxLevel=%d, scaleFactor=%.6f)",
			level, levelWidth, levelHeight, maxLevel, scaleFactor)
	}

	// Вычисляем максимальное количество тайлов
	maxCol := int(math.Max(1, math.Ceil(float64(levelWidth)/float64(s.tileSize))))
	maxRow := int(math.Max(1, math.Ceil(float64(levelHeight)/float64(s.tileSize))))

	// Проверяем границы
	if col < 0 || row < 0 || col >= maxCol || row >= maxRow {
		return nil, fmt.Errorf("tile coordinates out of bounds: level=%d, col=%d, row=%d (max: col=%d, row=%d, image_size=%dx%d)",
			level, col, row, maxCol-1, maxRow-1, levelWidth, levelHeight)
	}

	// Проверяем кэш тайлов
	tileCacheKey := fmt.Sprintf("%s_%d_%d_%d.%s", strings.ReplaceAll(imagePath, "/", "_"), level, col, row, format)
	tileCachePath := filepath.Join(s.tileCacheDir, tileCacheKey)

	if cachedTile, err := os.ReadFile(tileCachePath); err == nil {
		return &domain.Tile{
			Level:  level,
			Col:    col,
			Row:    row,
			Data:   cachedTile,
			Format: format,
		}, nil
	}

	// Загружаем изображение
	osImg, err := s.getOrLoadImage(ctx, imagePath)
	if err != nil {
		return nil, err
	}
	defer s.closeImage(osImg)

	// Вычисляем координаты тайла
	tileX := col * s.tileSize
	tileY := row * s.tileSize

	// Определяем, какой уровень OpenSlide использовать
	// OpenSlide использует встроенные уровни, которые могут не совпадать с DZI уровнями
	openSlideLevel := s.findBestOpenSlideLevel(osImg.osr, level, info)

	// Получаем размеры выбранного уровня OpenSlide
	var levelWidthOS, levelHeightOS C.int64_t
	C.openslide_get_level_dimensions(osImg.osr, C.int(openSlideLevel), &levelWidthOS, &levelHeightOS)

	// Вычисляем координаты в исходном изображении (уровень maxLevel = полное разрешение)
	// В OpenSeadragon: maxLevel = полное разрешение, 0 = минимальный масштаб
	// Формула для координат: sourceCoord = tileCoord / scaleFactor, где scaleFactor = 2^(level - maxLevel)
	maxLevelForCoords := info.Levels - 1
	coordScaleFactor := math.Pow(2, float64(level-maxLevelForCoords))
	if coordScaleFactor <= 0 {
		return nil, fmt.Errorf("invalid scale factor: level=%d, maxLevel=%d, scaleFactor=%.6f",
			level, maxLevelForCoords, coordScaleFactor)
	}

	sourceX := int(float64(tileX-s.overlap) / coordScaleFactor)
	sourceY := int(float64(tileY-s.overlap) / coordScaleFactor)
	sourceWidth := int(float64(s.tileSize+2*s.overlap) / coordScaleFactor)
	sourceHeight := int(float64(s.tileSize+2*s.overlap) / coordScaleFactor)

	// Проверяем валидность размеров источника
	if sourceWidth <= 0 || sourceHeight <= 0 {
		return nil, fmt.Errorf("invalid source dimensions: sourceWidth=%d, sourceHeight=%d (tileSize=%d, overlap=%d, coordScaleFactor=%.6f)",
			sourceWidth, sourceHeight, s.tileSize, s.overlap, coordScaleFactor)
	}

	// Ограничиваем координаты
	sourceX = int(math.Max(0, float64(sourceX)))
	sourceY = int(math.Max(0, float64(sourceY)))
	sourceWidth = int(math.Min(float64(info.Width-sourceX), float64(sourceWidth)))
	sourceHeight = int(math.Min(float64(info.Height-sourceY), float64(sourceHeight)))

	// Вычисляем масштаб между уровнем 0 и выбранным уровнем OpenSlide
	var level0Width, level0Height C.int64_t
	C.openslide_get_level0_dimensions(osImg.osr, &level0Width, &level0Height)

	levelScale := float64(level0Width) / float64(levelWidthOS)

	// Координаты на выбранном уровне OpenSlide
	levelX := int(float64(sourceX) / levelScale)
	levelY := int(float64(sourceY) / levelScale)
	levelW := int(float64(sourceWidth) / levelScale)
	levelH := int(float64(sourceHeight) / levelScale)

	// Проверяем, что размеры валидны
	if levelW <= 0 || levelH <= 0 {
		return nil, fmt.Errorf("invalid region size: levelW=%d, levelH=%d (sourceWidth=%d, sourceHeight=%d, levelScale=%.2f)",
			levelW, levelH, sourceWidth, sourceHeight, levelScale)
	}

	// Читаем область из OpenSlide
	// OpenSlide использует ARGB формат (32 бита на пиксель)
	// openslide_read_region(osr, dest, x, y, level, w, h)
	bufSize := levelW * levelH
	if bufSize <= 0 {
		return nil, fmt.Errorf("invalid buffer size: levelW=%d * levelH=%d = %d", levelW, levelH, bufSize)
	}

	buf := make([]uint32, bufSize)
	if len(buf) == 0 {
		return nil, fmt.Errorf("failed to allocate buffer: size=%d", bufSize)
	}

	C.openslide_read_region(osImg.osr, (*C.uint32_t)(unsafe.Pointer(&buf[0])), C.int64_t(levelX), C.int64_t(levelY), C.int(openSlideLevel), C.int64_t(levelW), C.int64_t(levelH))

	// Проверяем ошибки OpenSlide
	errStr := C.openslide_get_error(osImg.osr)
	if errStr != nil {
		return nil, fmt.Errorf("openslide error: %s", C.GoString(errStr))
	}

	// Конвертируем ARGB в image.Image
	img := s.argbToImage(buf, levelW, levelH)

	// Масштабируем до нужного размера тайла, если необходимо
	targetWidth := s.tileSize + 2*s.overlap
	targetHeight := s.tileSize + 2*s.overlap

	if img.Bounds().Dx() != targetWidth || img.Bounds().Dy() != targetHeight {
		img = s.resizeImage(img, targetWidth, targetHeight)
	}

	// Кодируем в нужный формат
	var tileData []byte
	var encodeErr error
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		tileData, encodeErr = s.encodeJPEG(img)
	case "png":
		tileData, encodeErr = s.encodePNG(img)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if encodeErr != nil {
		return nil, fmt.Errorf("failed to encode tile: %w", encodeErr)
	}

	// Сохраняем в кэш (асинхронно)
	go func() {
		if err := os.WriteFile(tileCachePath, tileData, 0644); err != nil {
			slog.Debug("failed to cache tile", "path", tileCachePath, "err", err)
		}
	}()

	return &domain.Tile{
		Level:  level,
		Col:    col,
		Row:    row,
		Data:   tileData,
		Format: format,
	}, nil
}

func (s *openslideService) GetDZI(ctx context.Context, imagePath string) (*domain.DZI, error) {
	info, err := s.GetImageInfo(ctx, imagePath)
	if err != nil {
		return nil, err
	}

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

func (s *openslideService) getOrLoadImage(ctx context.Context, imagePath string) (*openslideImage, error) {
	// Проверяем кэш
	s.cacheMux.RLock()
	cached, exists := s.cache[imagePath]
	s.cacheMux.RUnlock()

	if exists {
		cached.accessMux.Lock()
		cached.lastUsed = time.Now()
		cached.accessMux.Unlock()
		return cached, nil
	}

	// Загружаем файл из S3
	localPath := filepath.Join(s.cacheDir, strings.ReplaceAll(imagePath, "/", "_"))

	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		stream, err := s.s3Client.GetObjectStream(ctx, "", imagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to get image stream from S3: %w", err)
		}
		defer stream.Close()

		outFile, err := os.Create(localPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create cache file: %w", err)
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, stream); err != nil {
			os.Remove(localPath)
			return nil, fmt.Errorf("failed to save image to cache: %w", err)
		}
	}

	// Открываем через OpenSlide
	cPath := C.CString(localPath)
	defer C.free(unsafe.Pointer(cPath))

	osr := C.openslide_open(cPath)
	if osr == nil {
		return nil, fmt.Errorf("failed to open image with OpenSlide: %s", localPath)
	}

	// Проверяем ошибки
	errStr := C.openslide_get_error(osr)
	if errStr != nil {
		C.openslide_close(osr)
		return nil, fmt.Errorf("openslide error: %s", C.GoString(errStr))
	}

	osImg := &openslideImage{
		osr:      osr,
		path:     localPath,
		lastUsed: time.Now(),
	}

	// Добавляем в кэш
	s.cacheMux.Lock()
	if len(s.cache) >= s.maxCacheSize {
		s.cleanupCache()
	}
	s.cache[imagePath] = osImg
	s.cacheMux.Unlock()

	return osImg, nil
}

func (s *openslideService) closeImage(osImg *openslideImage) {
	if osImg != nil && osImg.osr != nil {
		// Не закрываем здесь, так как изображение может быть в кэше
		// Закрытие происходит в cleanupCache
	}
}

func (s *openslideService) cleanupCache() {
	now := time.Now()
	for path, cached := range s.cache {
		if now.Sub(cached.lastUsed) > s.cacheTTL {
			cached.accessMux.Lock()
			if cached.osr != nil {
				C.openslide_close(cached.osr)
				cached.osr = nil
			}
			cached.accessMux.Unlock()
			delete(s.cache, path)
		}
	}
}

func (s *openslideService) findBestOpenSlideLevel(osr *C.openslide_t, dziLevel int, info *domain.ImageInfo) int {
	levelCount := int(C.openslide_get_level_count(osr))

	// Вычисляем желаемый масштаб для DZI уровня
	desiredScale := math.Pow(2, float64(dziLevel))

	bestLevel := 0
	bestDiff := math.MaxFloat64

	// Ищем ближайший уровень OpenSlide
	for i := 0; i < levelCount; i++ {
		var levelWidth, levelHeight C.int64_t
		C.openslide_get_level_dimensions(osr, C.int(i), &levelWidth, &levelHeight)

		var level0Width, level0Height C.int64_t
		C.openslide_get_level0_dimensions(osr, &level0Width, &level0Height)

		levelScale := float64(level0Width) / float64(levelWidth)
		diff := math.Abs(levelScale - desiredScale)

		if diff < bestDiff {
			bestDiff = diff
			bestLevel = i
		}
	}

	return bestLevel
}

func (s *openslideService) argbToImage(argb []uint32, width, height int) image.Image {
	// Проверяем валидность размеров
	if width <= 0 || height <= 0 {
		// Возвращаем пустое изображение минимального размера
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}

	if len(argb) == 0 {
		// Возвращаем пустое изображение нужного размера
		return image.NewRGBA(image.Rect(0, 0, width, height))
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	expectedPixLen := width * height * 4
	if len(img.Pix) != expectedPixLen {
		// Если размер не совпадает, создаем новый с правильным размером
		img = image.NewRGBA(image.Rect(0, 0, width, height))
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x
			if idx < len(argb) {
				// ARGB (32-bit) -> RGBA
				argbVal := argb[idx]
				a := uint8((argbVal >> 24) & 0xFF)
				r := uint8((argbVal >> 16) & 0xFF)
				g := uint8((argbVal >> 8) & 0xFF)
				b := uint8(argbVal & 0xFF)

				offset := (y*width + x) * 4
				if offset+3 < len(img.Pix) {
					img.Pix[offset] = r
					img.Pix[offset+1] = g
					img.Pix[offset+2] = b
					img.Pix[offset+3] = a
				}
			}
		}
	}

	return img
}

func (s *openslideService) resizeImage(img image.Image, width, height int) image.Image {
	// Используем простой алгоритм масштабирования
	// Можно заменить на более качественный
	return imaging.Resize(img, width, height, imaging.Lanczos)
}

func (s *openslideService) encodeJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	return buf.Bytes(), err
}

func (s *openslideService) encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	return buf.Bytes(), err
}
