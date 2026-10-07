package original_image

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cytology/internal/domain"
	"cytology/internal/repository"
	repoentity "cytology/internal/repository/entity"
	original_imageEntity "cytology/internal/repository/original_image/entity"

	"github.com/google/uuid"
)

type Service interface {
	CreateOriginalImage(ctx context.Context, arg CreateOriginalImageArg) (uuid.UUID, error)
	GetOriginalImageByID(ctx context.Context, id uuid.UUID) (domain.OriginalImage, error)
	GetOriginalImagesByCytologyID(ctx context.Context, cytologyID uuid.UUID) ([]domain.OriginalImage, error)
	UpdateOriginalImage(ctx context.Context, arg UpdateOriginalImageArg) (domain.OriginalImage, error)
}

type service struct {
	dao repository.DAO
}

func New(dao repository.DAO) Service {
	return &service{
		dao: dao,
	}
}

type CreateOriginalImageArg struct {
	CytologyID  uuid.UUID
	File        []byte // Используется только если ImagePath не указан
	ContentType string
	DelayTime   *float64
	ImagePath   *string // Путь к файлу в S3 (если файл уже загружен)
}

type UpdateOriginalImageArg struct {
	Id         uuid.UUID
	DelayTime  *float64
	ViewedFlag *bool
}

func (s *service) CreateOriginalImage(ctx context.Context, arg CreateOriginalImageArg) (uuid.UUID, error) {
	var imageID uuid.UUID
	var imagePath string

	// Если передан путь к файлу, используем его (файл уже загружен в S3)
	if arg.ImagePath != nil && *arg.ImagePath != "" {
		imagePath = *arg.ImagePath
		// Извлекаем imageID из пути: {cytology_id}/{image_id}/{image_id}
		// Путь имеет формат: {cytology_id}/{image_id}/{image_id}
		parts := strings.Split(imagePath, "/")
		if len(parts) >= 2 {
			var err error
			imageID, err = uuid.Parse(parts[1])
			if err != nil {
				return uuid.Nil, fmt.Errorf("invalid image_id in image_path: %w", err)
			}
		} else {
			return uuid.Nil, fmt.Errorf("invalid image_path format: %s", imagePath)
		}
	} else {
		// Если путь не передан, генерируем ID и загружаем файл в S3
		imageID = uuid.New()

		// Формируем путь в S3: {cytology_id}/{image_id}/{image_id}
		// Используем "/" для S3, так как filepath.Join может давать разные результаты на разных ОС
		imagePath = arg.CytologyID.String() + "/" + imageID.String() + "/" + imageID.String()

		// Загружаем файл в S3
		fileRepo := s.dao.NewFileRepo()
		file := domain.File{
			Format: arg.ContentType,
			Size:   int64(len(arg.File)),
			Buf:    bytes.NewReader(arg.File),
		}

		if err := fileRepo.LoadFile(ctx, imagePath, file); err != nil {
			return uuid.Nil, fmt.Errorf("load file to S3: %w", err)
		}
	}

	// Создаем запись в БД
	img := domain.OriginalImage{
		Id:         imageID,
		CytologyID: arg.CytologyID,
		ImagePath:  imagePath,
		CreateDate: time.Now(),
		DelayTime:  arg.DelayTime,
		ViewedFlag: false,
	}

	entityImg := original_imageEntity.OriginalImage{}.FromDomain(img)
	if err := s.dao.NewOriginalImageQuery(ctx).InsertOriginalImage(entityImg); err != nil {
		return uuid.Nil, err
	}

	return img.Id, nil
}

func (s *service) GetOriginalImageByID(ctx context.Context, id uuid.UUID) (domain.OriginalImage, error) {
	img, err := s.dao.NewOriginalImageQuery(ctx).GetOriginalImageByID(id)
	if err != nil {
		return domain.OriginalImage{}, err
	}
	return img.ToDomain(), nil
}

func (s *service) GetOriginalImagesByCytologyID(ctx context.Context, cytologyID uuid.UUID) ([]domain.OriginalImage, error) {
	visited := make(map[uuid.UUID]struct{})
	for {
		if _, ok := visited[cytologyID]; ok {
			return nil, fmt.Errorf("cycle in cytology image history: %s", cytologyID)
		}
		visited[cytologyID] = struct{}{}

		images, err := s.dao.NewOriginalImageQuery(ctx).GetOriginalImagesByCytologyID(cytologyID)
		if err != nil && !errors.Is(err, repoentity.ErrNotFound) {
			return nil, err
		}
		if len(images) > 0 {
			return original_imageEntity.OriginalImage{}.SliceToDomain(images), nil
		}

		img, err := s.dao.NewCytologyImageQuery(ctx).GetCytologyImageByID(cytologyID)
		if err != nil {
			return nil, err
		}
		if !img.PrevID.Valid {
			return nil, repoentity.ErrNotFound
		}
		cytologyID = img.PrevID.UUID
	}
}

func (s *service) UpdateOriginalImage(ctx context.Context, arg UpdateOriginalImageArg) (domain.OriginalImage, error) {
	img, err := s.dao.NewOriginalImageQuery(ctx).GetOriginalImageByID(arg.Id)
	if err != nil {
		return domain.OriginalImage{}, err
	}

	domainImg := img.ToDomain()
	if arg.DelayTime != nil {
		domainImg.DelayTime = arg.DelayTime
	}
	if arg.ViewedFlag != nil {
		domainImg.ViewedFlag = *arg.ViewedFlag
	}

	entityImg := original_imageEntity.OriginalImage{}.FromDomain(domainImg)
	if err := s.dao.NewOriginalImageQuery(ctx).UpdateOriginalImage(entityImg); err != nil {
		return domain.OriginalImage{}, err
	}

	return domainImg, nil
}
