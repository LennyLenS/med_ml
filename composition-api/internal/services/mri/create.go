package mri

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"

	adapter "composition-api/internal/adapters/exam"
	mriuploadpb "composition-api/internal/generated/dbus/produce/mriupload"
)

func (s *service) Create(ctx context.Context, in CreateMriArg) (uuid.UUID, error) {
	mriID, err := s.adapters.Exam.CreateMri(ctx, adapter.CreateMriIn{
		Projection:  in.Projection,
		ExternalID:  in.ExternalID,
		Author:      in.Author,
		DeviceID:    in.DeviceID,
		Description: in.Description,
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("create mri in microservice: %w", err)
	}

	err = s.dao.NewFileRepo().LoadFile(ctx, filepath.Join(mriID.String(), mriID.String()), in.File)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("load mri file to s3: %w", err)
	}

	// TODO: сделать сагу
	err = s.dbus.SendMriUpload(ctx, &mriuploadpb.MriUpload{MriId: mriID.String()})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("send mri upload to dbus: %w", err)
	}

	return mriID, nil
}
