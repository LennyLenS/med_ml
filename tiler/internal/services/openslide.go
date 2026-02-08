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
		return nil, fmt.Errorf("failed to load image for info (path: %s): %w", imagePath, err)
	}

	// Защищаем доступ к osr мьютексом (OpenSlide не thread-safe)
	osImg.accessMux.RLock()

	// Проверяем, что osr не был закрыт
	if osImg.osr == nil {
		osImg.accessMux.RUnlock()
		return nil, fmt.Errorf("image handle was closed")
	}

	// Сохраняем указатель локально для безопасности
	osr := osImg.osr

	defer osImg.accessMux.RUnlock()

	// Получаем размеры уровня 0 (полное разрешение)
	var width, height C.int64_t
	C.openslide_get_level0_dimensions(osr, &width, &height)

	// Получаем количество уровней
	levelCount := int(C.openslide_get_level_count(osr))

	// Вычисляем количество уровней для DZI согласно OpenSeadragon
	// OpenSeadragon генерирует уровни до размера 1x1
	// Формула: находим maxLevel такой, что maxDim / 2^maxLevel >= 1
	// Т.е. 2^maxLevel <= maxDim, значит maxLevel = floor(log2(maxDim))
	maxDim := math.Max(float64(width), float64(height))
	log2MaxDim := math.Log2(maxDim)
	maxLevel := int(math.Floor(log2MaxDim))

	// Проверяем, что на level 0 изображение >= 1px
	powerMaxLevel := math.Pow(2, float64(maxLevel))
	if powerMaxLevel > maxDim {
		// 2^maxLevel > maxDim, значит нужно уменьшить maxLevel
		maxLevel--
	}

	// Количество уровней = maxLevel + 1 (от 0 до maxLevel включительно)
	levels := maxLevel + 1

	// Используем максимум из встроенных уровней OpenSlide и вычисленных для DZI
	// Но не меньше, чем нужно для достижения 1x1
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
	// Защита от паники
	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetTile: panic recovered", "err", r, "imagePath", imagePath, "level", level, "col", col, "row", row)
		}
	}()

	info, err := s.GetImageInfo(ctx, imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get image info: %w", err)
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

	// OpenSeadragon требует минимум 1x1 пиксель для любого уровня
	// Если из-за округления получился 0, устанавливаем 1
	if levelWidth <= 0 {
		levelWidth = 1
	}
	if levelHeight <= 0 {
		levelHeight = 1
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
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	if osImg == nil {
		return nil, fmt.Errorf("invalid image handle: osImg is nil")
	}

	// Обновляем время последнего использования
	osImg.lastUsed = time.Now()

	// Защищаем доступ к osr мьютексом (OpenSlide не thread-safe)
	osImg.accessMux.RLock()

	// Проверяем, что osr не был закрыт (после получения блокировки)
	if osImg.osr == nil {
		osImg.accessMux.RUnlock()
		return nil, fmt.Errorf("image handle was closed")
	}

	// Сохраняем указатель локально для безопасности
	osr := osImg.osr

	defer osImg.accessMux.RUnlock()

	// Вычисляем координаты тайла
	tileX := col * s.tileSize
	tileY := row * s.tileSize

	// Определяем, какой уровень OpenSlide использовать
	// OpenSlide использует встроенные уровни, которые могут не совпадать с DZI уровнями
	openSlideLevel := s.findBestOpenSlideLevel(osr, level, info)

	// Проверяем валидность уровня OpenSlide
	levelCount := int(C.openslide_get_level_count(osr))
	if openSlideLevel < 0 || openSlideLevel >= levelCount {
		return nil, fmt.Errorf("invalid OpenSlide level: %d (max: %d, dziLevel: %d)", openSlideLevel, levelCount-1, level)
	}

	// Получаем размеры выбранного уровня OpenSlide
	var levelWidthOS, levelHeightOS C.int64_t
	C.openslide_get_level_dimensions(osr, C.int(openSlideLevel), &levelWidthOS, &levelHeightOS)

	// Проверяем ошибки OpenSlide после получения размеров
	if errStr := C.openslide_get_error(osr); errStr != nil {
		return nil, fmt.Errorf("openslide error after get_level_dimensions: %s", C.GoString(errStr))
	}

	// Вычисляем масштаб между уровнем 0 и выбранным уровнем OpenSlide
	var level0Width, level0Height C.int64_t
	C.openslide_get_level0_dimensions(osr, &level0Width, &level0Height)

	// Проверяем на деление на ноль
	if levelWidthOS <= 0 || levelHeightOS <= 0 {
		return nil, fmt.Errorf("invalid OpenSlide level dimensions: levelWidth=%d, levelHeight=%d", levelWidthOS, levelHeightOS)
	}

	levelScale := float64(level0Width) / float64(levelWidthOS)
	if levelScale <= 0 {
		return nil, fmt.Errorf("invalid level scale: level0Width=%d, levelWidthOS=%d, scale=%.6f", level0Width, levelWidthOS, levelScale)
	}

	// Вычисляем координаты в исходном изображении (уровень maxLevel = полное разрешение)
	// В OpenSeadragon: maxLevel = полное разрешение, 0 = минимальный масштаб
	// Формула для координат: sourceCoord = (tileCoord - overlap) / scaleFactor
	// ВАЖНО: overlap ВЫЧИТАЕТСЯ из координат, чтобы тайлы правильно перекрывались
	maxLevelForCoords := info.Levels - 1
	coordScaleFactor := math.Pow(2, float64(level-maxLevelForCoords))
	if coordScaleFactor <= 0 {
		return nil, fmt.Errorf("invalid scale factor: level=%d, maxLevel=%d, scaleFactor=%.6f",
			level, maxLevelForCoords, coordScaleFactor)
	}

	// Координаты тайла С вычитанием overlap (для правильного перекрытия тайлов)
	sourceX := int(float64(tileX-s.overlap) / coordScaleFactor)
	sourceY := int(float64(tileY-s.overlap) / coordScaleFactor)
	// Размер области С overlap (overlap добавляется к размеру тайла)
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

	// Координаты на выбранном уровне OpenSlide
	// Масштабируем координаты из полного разрешения на выбранный уровень OpenSlide
	levelX := int(float64(sourceX) / levelScale)
	levelY := int(float64(sourceY) / levelScale)
	levelW := int(float64(sourceWidth) / levelScale)
	levelH := int(float64(sourceHeight) / levelScale)

	// Для низких уровней DZI размер области может быть очень большим
	// Ограничиваем размер области размером тайла на выбранном уровне OpenSlide
	// Это предотвращает чтение огромных областей для низких уровней
	targetTileSize := s.tileSize + 2*s.overlap
	if levelW > targetTileSize*2 {
		levelW = targetTileSize * 2
	}
	if levelH > targetTileSize*2 {
		levelH = targetTileSize * 2
	}

	// Проверяем, что координаты не отрицательные
	if levelX < 0 {
		levelX = 0
	}
	if levelY < 0 {
		levelY = 0
	}

	// Проверяем, что размеры валидны
	if levelW <= 0 || levelH <= 0 {
		return nil, fmt.Errorf("invalid region size: levelW=%d, levelH=%d (sourceWidth=%d, sourceHeight=%d, levelScale=%.6f, levelX=%d, levelY=%d)",
			levelW, levelH, sourceWidth, sourceHeight, levelScale, levelX, levelY)
	}

	// Проверяем, что координаты не выходят за границы уровня OpenSlide
	if int64(levelX)+int64(levelW) > int64(levelWidthOS) {
		levelW = int(levelWidthOS) - levelX
		if levelW <= 0 {
			return nil, fmt.Errorf("region X coordinate out of bounds: levelX=%d, levelW=%d, levelWidthOS=%d", levelX, levelW, levelWidthOS)
		}
	}
	if int64(levelY)+int64(levelH) > int64(levelHeightOS) {
		levelH = int(levelHeightOS) - levelY
		if levelH <= 0 {
			return nil, fmt.Errorf("region Y coordinate out of bounds: levelY=%d, levelH=%d, levelHeightOS=%d", levelY, levelH, levelHeightOS)
		}
	}

	// Читаем область из OpenSlide
	// OpenSlide использует ARGB формат (32 бита на пиксель)
	// openslide_read_region(osr, dest, x, y, level, w, h)
	bufSize := levelW * levelH
	if bufSize <= 0 {
		return nil, fmt.Errorf("invalid buffer size: levelW=%d * levelH=%d = %d", levelW, levelH, bufSize)
	}

	// Проверяем максимальный размер буфера (защита от переполнения памяти)
	const maxBufferSize = 100 * 1024 * 1024 // 100MB (примерно 10000x10000 пикселей)
	if bufSize > maxBufferSize {
		return nil, fmt.Errorf("buffer size too large: %d (max: %d, levelW=%d, levelH=%d)", bufSize, maxBufferSize, levelW, levelH)
	}

	buf := make([]uint32, bufSize)
	if len(buf) == 0 {
		return nil, fmt.Errorf("failed to allocate buffer: size=%d", bufSize)
	}

	// Проверяем, что osr все еще валиден перед вызовом CGO
	if osr == nil {
		return nil, fmt.Errorf("osr became nil before read_region")
	}

	C.openslide_read_region(osr, (*C.uint32_t)(unsafe.Pointer(&buf[0])), C.int64_t(levelX), C.int64_t(levelY), C.int(openSlideLevel), C.int64_t(levelW), C.int64_t(levelH))

	// Проверяем ошибки OpenSlide
	if errStr := C.openslide_get_error(osr); errStr != nil {
		return nil, fmt.Errorf("openslide error after read_region: %s", C.GoString(errStr))
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
		return nil, fmt.Errorf("failed to get image info for DZI (path: %s): %w", imagePath, err)
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
		// Создаем директорию для кэша, если её нет
		if err := os.MkdirAll(s.cacheDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create cache directory: %w", err)
		}

		stream, err := s.s3Client.GetObjectStream(ctx, "", imagePath)
		if err != nil {
			// Проверяем, является ли ошибка "key does not exist"
			if strings.Contains(err.Error(), "does not exist") ||
				strings.Contains(err.Error(), "NoSuchKey") ||
				strings.Contains(err.Error(), "not found") {
				return nil, fmt.Errorf("image not found in S3: %s (path: %s)", imagePath, imagePath)
			}
			return nil, fmt.Errorf("failed to get image stream from S3 (path: %s): %w", imagePath, err)
		}
		defer stream.Close()

		outFile, err := os.Create(localPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create cache file: %w", err)
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, stream); err != nil {
			os.Remove(localPath)
			return nil, fmt.Errorf("failed to save image to cache (localPath: %s): %w", localPath, err)
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
	s.cacheMux.Lock()
	defer s.cacheMux.Unlock()

	for path, cached := range s.cache {
		if now.Sub(cached.lastUsed) > s.cacheTTL {
			// Получаем эксклюзивную блокировку перед закрытием
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
	if osr == nil {
		slog.Error("findBestOpenSlideLevel: osr is nil")
		return 0
	}

	levelCount := int(C.openslide_get_level_count(osr))
	if levelCount <= 0 {
		slog.Error("findBestOpenSlideLevel: invalid level count", "count", levelCount)
		return 0
	}

	// Вычисляем желаемый масштаб для DZI уровня
	// В OpenSeadragon: maxLevel = полное разрешение, 0 = минимальный масштаб
	// Формула: desiredScale = 2^(dziLevel - maxLevel)
	maxLevel := info.Levels - 1
	desiredScale := math.Pow(2, float64(dziLevel-maxLevel))

	bestLevel := 0
	bestDiff := math.MaxFloat64

	// Ищем ближайший уровень OpenSlide
	for i := 0; i < levelCount; i++ {
		var levelWidth, levelHeight C.int64_t
		C.openslide_get_level_dimensions(osr, C.int(i), &levelWidth, &levelHeight)

		// Проверяем ошибки OpenSlide
		if errStr := C.openslide_get_error(osr); errStr != nil {
			slog.Warn("findBestOpenSlideLevel: openslide error", "level", i, "error", C.GoString(errStr))
			continue
		}

		var level0Width, level0Height C.int64_t
		C.openslide_get_level0_dimensions(osr, &level0Width, &level0Height)

		// Проверяем на деление на ноль
		if levelWidth <= 0 || levelHeight <= 0 {
			slog.Warn("findBestOpenSlideLevel: invalid level dimensions", "level", i, "width", levelWidth, "height", levelHeight)
			continue
		}

		levelScale := float64(level0Width) / float64(levelWidth)
		if levelScale <= 0 {
			slog.Warn("findBestOpenSlideLevel: invalid level scale", "level", i, "scale", levelScale)
			continue
		}

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
