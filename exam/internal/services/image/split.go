package image

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	mrisplittedpb "exam/internal/generated/dbus/produce/mrisplitted"

	"exam/internal/domain"
	"exam/internal/repository/image/entity"
	"exam/internal/services/splitter"

	"github.com/google/uuid"
)

var ErrMriDoesntExist = errors.New("mri doesnt exist")

func (s *service) SplitMri(ctx context.Context, id uuid.UUID) error {
	fileRepo := s.dao.NewFileRepo()

	exists, err := s.dao.NewMriQuery(ctx).CheckExist(id)
	if err != nil {
		return fmt.Errorf("check exists mri: %w", err)
	}
	if !exists {
		return ErrMriDoesntExist
	}

	file, closer, err := fileRepo.GetFileViaTemp(ctx, filepath.Join(id.String(), id.String()))
	if err != nil {
		return fmt.Errorf("get file via temp: %w", err)
	}
	defer closer()

	splitterSrv := splitter.New()
	splitted, err := splitterSrv.SplitToPng(file)
	if err != nil {
		return fmt.Errorf("split img to png: %w", err)
	}

	images := make([]domain.Image, len(splitted))
	for i := range images {
		images[i].Id = uuid.New()
		images[i].MriID = id
		images[i].Page = i + 1
	}

	for i, image := range images {
		if err := fileRepo.LoadFile(
			ctx,
			filepath.Join(id.String(), image.Id.String(), image.Id.String()),
			splitted[i],
		); err != nil {
			return fmt.Errorf("load file to S3: %w", err)
		}
	}

	ctx, err = s.dao.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = s.dao.RollbackTx(ctx) }()

	if err := s.dao.NewImageQuery(ctx).InsertImages(entity.Image{}.SliceFromDomain(images)...); err != nil {
		return fmt.Errorf("insert images: %w", err)
	}

	if err := s.dao.NewMriQuery(ctx).UpdateMriStatus(id, string(domain.ExamStatusPending)); err != nil {
		return fmt.Errorf("update mri status: %w", err)
	}

	if err := s.dao.CommitTx(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	imageIds := make([]uuid.UUID, 0, len(images))
	for _, image := range images {
		imageIds = append(imageIds, image.Id)
	}

	if err := s.dbus.SendMriSplitted(ctx, &mrisplittedpb.MriSplitted{
		MriId:   id.String(),
		PagesId: uuid.UUIDs(imageIds).Strings(),
	}); err != nil {
		return fmt.Errorf("send to mrisplitted topic: %w", err)
	}

	return nil
}
