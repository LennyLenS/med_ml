package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"log/slog"
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
	infoCache    map[string]*domain.ImageInfo // Кэш информации об изображениях
	infoCacheMux sync.RWMutex
	cacheDir     string
	tileCacheDir string // Директория для кэширования тайлов
	maxCacheSize int
	cacheTTL     time.Duration
}

type S3Client interface {
	GetObject(ctx context.Context, bucketName, objectName string) ([]byte, error)
	GetObjectStream(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
}

func NewImageService(s3Client S3Client, tileSize, overlap int) ImageService {
	// ВРЕМЕННО ЗАКОММЕНТИРОВАНО для тестирования OpenSlide
	// Инициализируем libvips (требуется только один раз)
	// vips.Startup(nil)

	// Создаем временную директорию для кэша файлов
	cacheDir := filepath.Join(os.TempDir(), "tiler_cache")
	os.MkdirAll(cacheDir, 0755)

	// Создаем директорию для кэширования тайлов
	tileCacheDir := filepath.Join(os.TempDir(), "tiler_tile_cache")
	os.MkdirAll(tileCacheDir, 0755)

	return &imageService{
		s3Client:     s3Client,
		tileSize:     tileSize,
		overlap:      overlap,
		cache:        make(map[string]*cachedImage),
		infoCache:    make(map[string]*domain.ImageInfo),
		cacheDir:     cacheDir,
		tileCacheDir: tileCacheDir,
		maxCacheSize: 10,        // Максимум 10 открытых файлов в кэше
		cacheTTL:     time.Hour, // Время жизни кэша
	}
}

// GetImageInfo возвращает информацию об изображении, включая максимальный уровень масштабирования
//
// ФОРМУЛА OPENSEADRAGON DZI:
// Используется самая длинная сторона изображения (M).
// Формула: 2^(N+1) ≥ M, где M — длина самой длинной стороны.
// Для изображения 197208 px:
//
//	2^17 = 131072 — мало
//	2^18 = 262144 — подходит
//	→ maxLevel = 17 (N = 17)
//	→ количество уровней = 18 (от 0 до 17)
//
// КАК РАБОТАЮТ УРОВНИ:
// - Уровень maxLevel (17) — полное разрешение (197208 × 88437 px)
// - Уровень maxLevel-1 (16) — ½ (около 98604 × 44219 px)
// - Уровень maxLevel-2 (15) — ¼ и т.д.
// - Уровень 0 — минимальный масштаб
// Каждый следующий уровень вдвое меньше предыдущего.
//
// ПАРАМЕТРЫ:
// - TileSize = 510 — каждый тайл 510×510 px
// - Overlap = 1 — соседние тайлы перекрываются на 1 px по краям
//
// ДЛЯ SVS ФАЙЛОВ:
// - SVS файлы (формат Aperio) уже содержат встроенные пирамидальные уровни
// - libvips автоматически использует эти уровни при чтении файла (random access)
// - Это позволяет эффективно читать только нужные тайлы без загрузки всего файла
//
// ПРИМЕРЫ РАСЧЕТА:
// - Изображение 197208×88437: maxLevel = 17, уровней = 18 (от 0 до 17)
// - Изображение 100000×80000: maxLevel = 16, уровней = 17 (от 0 до 16)
// - Изображение 10000×8000: maxLevel = 13, уровней = 14 (от 0 до 13)
//
// КАК УЗНАТЬ МАКСИМАЛЬНЫЙ УРОВЕНЬ:
// 1. Вызовите GetImageInfo() - вернет структуру ImageInfo с полем Levels
// 2. Максимальный доступный уровень = Levels - 1 (так как нумерация с 0)
// 3. Например, если Levels = 18, то доступны уровни от 0 до 17
func (s *imageService) GetImageInfo(ctx context.Context, imagePath string) (*domain.ImageInfo, error) {
	// Проверяем кэш информации об изображении
	s.infoCacheMux.RLock()
	cachedInfo, exists := s.infoCache[imagePath]
	s.infoCacheMux.RUnlock()

	if exists && cachedInfo != nil {
		// Возвращаем кэшированную информацию
		return &domain.ImageInfo{
			Width:    cachedInfo.Width,
			Height:   cachedInfo.Height,
			Levels:   cachedInfo.Levels,
			TileSize: cachedInfo.TileSize,
			Overlap:  cachedInfo.Overlap,
		}, nil
	}

	// Используем кэшированное изображение для получения метаданных
	// libvips читает только заголовок файла, не весь файл
	img, err := s.getOrLoadImage(ctx, imagePath)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	width := img.Width()
	height := img.Height()

	// Вычисляем количество уровней (levels) по формуле Deep Zoom Image (DZI)
	// Формула: levels = ceil(log2(max(width, height) / tileSize)) + 1
	// Каждый уровень масштабирования уменьшает изображение в 2 раза
	// Максимальный уровень - это когда изображение становится меньше или равно размеру тайла
	levels := s.calculateMaxLevels(width, height)

	// Для очень маленьких изображений гарантируем минимум 1 уровень
	if levels < 1 {
		levels = 1
	}

	info := &domain.ImageInfo{
		Width:    width,
		Height:   height,
		Levels:   levels,
		TileSize: s.tileSize,
		Overlap:  s.overlap,
	}

	// Сохраняем в кэш
	s.infoCacheMux.Lock()
	s.infoCache[imagePath] = info
	s.infoCacheMux.Unlock()

	// Возвращаем копию
	return &domain.ImageInfo{
		Width:    info.Width,
		Height:   info.Height,
		Levels:   info.Levels,
		TileSize: info.TileSize,
		Overlap:  info.Overlap,
	}, nil
}

// calculateMaxLevels вычисляет максимальное количество уровней масштабирования
// для изображения заданного размера согласно спецификации OpenSeadragon DZI
//
// Формула OpenSeadragon: 2^(N+1) ≥ M, где M — длина самой длинной стороны
// Для изображения 197208 px:
//   2^17 = 131072 — мало
//   2^18 = 262144 — подходит
//   → maxLevel = 17 (N = 17)
//   → количество уровней = 18 (от 0 до 17)
//
// Уровень 17 — полное разрешение (197208 × 88437 px)
// Уровень 16 — ½ (около 98604 × 44219 px)
// Уровень 15 — ¼ и т.д.
// Уровень 0 — минимальный масштаб
//
// Каждый следующий уровень вдвое меньше предыдущего.

func (s *imageService) calculateMaxLevels(width, height int) int {
	maxDim := math.Max(float64(width), float64(height))
	if maxDim <= 0 {
		return 1
	}

	// OpenSeadragon генерирует уровни до размера 1x1
	// Формула: находим maxLevel такой, что maxDim / 2^maxLevel >= 1
	// Т.е. 2^maxLevel <= maxDim, значит maxLevel = floor(log2(maxDim))
	// Это гарантирует, что на level 0 изображение будет >= 1px
	log2MaxDim := math.Log2(maxDim)
	maxLevel := int(math.Floor(log2MaxDim))

	// Проверяем, что на level 0 изображение >= 1px
	// Если нет, увеличиваем maxLevel на 1
	powerMaxLevel := math.Pow(2, float64(maxLevel))
	if powerMaxLevel > maxDim {
		// 2^maxLevel > maxDim, значит нужно уменьшить maxLevel
		maxLevel--
	}

	// Количество уровней = maxLevel + 1 (от 0 до maxLevel включительно)
	// Level maxLevel = полное разрешение
	// Level 0 = минимальный масштаб (>= 1px)
	levels := maxLevel + 1

	// Гарантируем минимум 1 уровень
	if levels < 1 {
		return 1
	}
	return levels
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
		// Используем streaming для загрузки из S3 (более эффективно для больших файлов)
		stream, err := s.s3Client.GetObjectStream(ctx, "", imagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to get image stream from S3 (path: %s): %w", imagePath, err)
		}
		defer stream.Close()

		// Создаем файл и копируем данные напрямую из потока
		outFile, err := os.Create(localPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create cache file: %w", err)
		}
		defer outFile.Close()

		// Копируем данные из потока в файл
		_, err = io.Copy(outFile, stream)
		if err != nil {
			os.Remove(localPath) // Удаляем частично загруженный файл
			return nil, fmt.Errorf("failed to save image to cache: %w", err)
		}
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
	// Включаем информацию об изображении в сообщение об ошибке для отладки
	if level < 0 || level >= info.Levels {
		return nil, fmt.Errorf("invalid level: %d (max: %d, image_size: %dx%d, calculated_levels: %d)",
			level, info.Levels-1, info.Width, info.Height, info.Levels)
	}

	// Вычисляем размеры изображения на данном уровне
	// В OpenSeadragon: maxLevel (например, 17) = полное разрешение
	//                  maxLevel - 1 (16) = уменьшено в 2 раза
	//                  maxLevel - 2 (15) = уменьшено в 4 раза
	//                  ...
	//                  0 = уменьшено в 2^maxLevel раз
	// Формула: scaleFactor = 2^(level - maxLevel)
	maxLevel := info.Levels - 1
	scaleFactor := math.Pow(2, float64(level-maxLevel))
	levelWidth := int(float64(info.Width) * scaleFactor)
	levelHeight := int(float64(info.Height) * scaleFactor)

	// Проверяем, что размеры изображения валидны
	if levelWidth <= 0 || levelHeight <= 0 {
		return nil, fmt.Errorf("invalid image dimensions at level %d: %dx%d", level, levelWidth, levelHeight)
	}

	// Вычисляем максимальное количество тайлов на данном уровне
	// Даже если изображение меньше тайла, должен быть хотя бы один тайл
	maxCol := int(math.Max(1, math.Ceil(float64(levelWidth)/float64(s.tileSize))))
	maxRow := int(math.Max(1, math.Ceil(float64(levelHeight)/float64(s.tileSize))))

	// Проверяем границы координат тайла
	// Если запрашивается тайл вне границ изображения, возвращаем ошибку
	if col < 0 || row < 0 || col >= maxCol || row >= maxRow {
		return nil, fmt.Errorf("tile coordinates out of bounds: level=%d, col=%d, row=%d (max: col=%d, row=%d, image_size=%dx%d)",
			level, col, row, maxCol-1, maxRow-1, levelWidth, levelHeight)
	}

	// Проверяем кэш тайлов на диске
	tileCacheKey := fmt.Sprintf("%s_%d_%d_%d.%s", strings.ReplaceAll(imagePath, "/", "_"), level, col, row, format)
	tileCachePath := filepath.Join(s.tileCacheDir, tileCacheKey)

	// Пытаемся прочитать тайл из кэша
	if cachedTile, err := os.ReadFile(tileCachePath); err == nil {
		return &domain.Tile{
			Level:  level,
			Col:    col,
			Row:    row,
			Data:   cachedTile,
			Format: format,
		}, nil
	}

	// ОПТИМИЗАЦИЯ: Извлекаем область из исходного файла ДО масштабирования
	// Это позволяет libvips использовать random access для tiled TIFF/SVS
	// и читать только нужные тайлы, а не загружать весь файл в память

	// Получаем изображение из кэша или загружаем его
	// libvips будет читать только нужные тайлы из файла благодаря random access
	baseImg, err := s.getOrLoadImage(ctx, imagePath)
	if err != nil {
		return nil, err
	}
	defer baseImg.Close()

	// Вычисляем координаты тайла на целевом уровне
	tileX := col * s.tileSize
	tileY := row * s.tileSize

	// Вычисляем масштаб для уровня
	// В OpenSeadragon: maxLevel = полное разрешение, 0 = минимальный масштаб
	// Формула: scaleFactor = 2^(level - maxLevel)
	maxLevelForTile := info.Levels - 1
	var tileScaleFactor float64
	if level == maxLevelForTile {
		tileScaleFactor = 1.0 // Полное разрешение
	} else {
		tileScaleFactor = math.Pow(2, float64(level-maxLevelForTile))
	}

	// Вычисляем координаты и размеры области в исходном изображении
	// С учетом overlap и масштаба
	sourceX := int(float64(tileX-s.overlap) / tileScaleFactor)
	sourceY := int(float64(tileY-s.overlap) / tileScaleFactor)
	sourceWidth := int(float64(s.tileSize+2*s.overlap) / tileScaleFactor)
	sourceHeight := int(float64(s.tileSize+2*s.overlap) / tileScaleFactor)

	// Ограничиваем координаты границами изображения
	sourceX = int(math.Max(0, float64(sourceX)))
	sourceY = int(math.Max(0, float64(sourceY)))
	sourceWidth = int(math.Min(float64(baseImg.Width()-sourceX), float64(sourceWidth)))
	sourceHeight = int(math.Min(float64(baseImg.Height()-sourceY), float64(sourceHeight)))

	// Проверяем, что область валидна
	if sourceX >= baseImg.Width() || sourceY >= baseImg.Height() || sourceWidth <= 0 || sourceHeight <= 0 {
		return nil, errors.New("tile coordinates out of bounds")
	}

	// Извлекаем область из исходного изображения ПЕРЕД масштабированием
	// Для tiled TIFF/SVS libvips использует random access и читает только нужные тайлы
	// ExtractArea модифицирует исходное изображение, поэтому нужно создать копию
	regionImg, err := baseImg.Copy()
	if err != nil {
		return nil, fmt.Errorf("failed to copy image: %w", err)
	}
	defer regionImg.Close()

	if err := regionImg.ExtractArea(sourceX, sourceY, sourceWidth, sourceHeight); err != nil {
		return nil, fmt.Errorf("failed to extract region: %w", err)
	}

	// Масштабируем только извлеченную область до размера тайла
	// Это намного эффективнее, чем масштабировать всё изображение
	// В OpenSeadragon: maxLevel = полное разрешение, 0 = минимальный масштаб
	var tileImg *vips.ImageRef
	if level == maxLevelForTile {
		// Для максимального уровня масштабирование не нужно
		tileImg, err = regionImg.Copy()
		if err != nil {
			return nil, fmt.Errorf("failed to copy region: %w", err)
		}
	} else {
		// Масштабируем только извлеченную область
		targetWidth := s.tileSize + 2*s.overlap
		targetHeight := s.tileSize + 2*s.overlap

		// Вычисляем масштаб для области
		regionScaleX := float64(targetWidth) / float64(sourceWidth)
		regionScaleY := float64(targetHeight) / float64(sourceHeight)
		regionScale := math.Min(regionScaleX, regionScaleY)

		// Используем более быстрый алгоритм для больших уменьшений
		var kernel vips.Kernel
		if regionScale < 0.25 {
			kernel = vips.KernelLinear // Быстрее для больших уменьшений
		} else {
			kernel = vips.KernelLanczos3 // Качественнее для небольших уменьшений
		}

		if err := regionImg.Resize(regionScale, kernel); err != nil {
			return nil, fmt.Errorf("failed to resize region: %w", err)
		}
		tileImg, err = regionImg.Copy()
		if err != nil {
			return nil, fmt.Errorf("failed to copy scaled region: %w", err)
		}
	}
	defer tileImg.Close()

	// Кодируем тайл в нужный формат
	var tileData []byte
	var encodeErr error
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		ep := vips.NewJpegExportParams()
		ep.Quality = 85 // Оптимизируем качество для баланса между размером и качеством
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

	// Сохраняем тайл в кэш на диске (асинхронно, чтобы не блокировать ответ)
	go func() {
		if err := os.WriteFile(tileCachePath, tileData, 0644); err != nil {
			// Игнорируем ошибки записи в кэш, чтобы не влиять на основной поток
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

// getTileFromStandardImage - fallback метод для стандартного image.Decode
func (s *imageService) getTileFromStandardImage(img image.Image, level, col, row int, format string, maxLevel int) (*domain.Tile, error) {
	// Масштабируем изображение до нужного уровня
	scaledImg := s.scaleToLevel(img, level, maxLevel)

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

func (s *imageService) scaleToLevel(img image.Image, level int, maxLevel int) image.Image {
	if level == maxLevel {
		return img
	}

	// Вычисляем масштаб для уровня
	// В OpenSeadragon: maxLevel = полное разрешение, 0 = минимальный масштаб
	// Формула: scaleFactor = 2^(level - maxLevel)
	scaleFactor := math.Pow(2, float64(level-maxLevel))
	bounds := img.Bounds()
	newWidth := int(float64(bounds.Dx()) * scaleFactor)
	newHeight := int(float64(bounds.Dy()) * scaleFactor)

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

// GetLevelInfo возвращает детальную информацию о конкретном уровне масштабирования
// Полезно для отладки и понимания структуры пирамиды изображения
func (s *imageService) GetLevelInfo(ctx context.Context, imagePath string, level int) (map[string]interface{}, error) {
	info, err := s.GetImageInfo(ctx, imagePath)
	if err != nil {
		return nil, err
	}

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

	// Вычисляем количество тайлов на данном уровне
	maxCol := int(math.Max(1, math.Ceil(float64(levelWidth)/float64(s.tileSize))))
	maxRow := int(math.Max(1, math.Ceil(float64(levelHeight)/float64(s.tileSize))))
	totalTiles := maxCol * maxRow

	return map[string]interface{}{
		"level":         level,
		"width":         levelWidth,
		"height":        levelHeight,
		"tiles_cols":    maxCol,
		"tiles_rows":    maxRow,
		"total_tiles":   totalTiles,
		"tile_size":     s.tileSize,
		"scale_factor":  scaleFactor,
		"is_max_level":  level == maxLevel,
		"original_size": fmt.Sprintf("%dx%d", info.Width, info.Height),
	}, nil
}
