package kt

import (
	"context"
	"fmt"
	"path/filepath"

	adapter "composition-api/internal/adapters/exam"

	ktuploadpb "composition-api/internal/generated/dbus/produce/ktupload"

	"github.com/google/uuid"
)

func (s *service) Create(ctx context.Context, in CreateKtArg) (uuid.UUID, error) {
	ktID, err := s.adapters.Exam.CreateKt(ctx, adapter.CreateKtIn{
		Author:      in.Author,
		DeviceID:    in.DeviceID,
		Description: in.Description,
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("create kt in microservice: %w", err)
	}

	err = s.dao.NewFileRepo().LoadFile(ctx, filepath.Join(ktID.String(), ktID.String()), in.File)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("load kt file to s3: %w", err)
	}

	err = s.dbus.SendKtUpload(ctx, &ktuploadpb.KtUpload{KtId: ktID.String()})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("send kt upload to dbus: %w", err)
	}

	return ktID, nil
}
