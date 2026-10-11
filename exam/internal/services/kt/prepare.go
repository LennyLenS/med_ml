package kt

import (
	"context"
	"fmt"
	"path/filepath"

	"exam/internal/domain"
	"exam/internal/generated/dbus/produce/ktprepared"
	"exam/internal/services/splitter"

	"github.com/google/uuid"
)

func (s *service) PrepareKt(ctx context.Context, id uuid.UUID) error {
	repo := s.dao.NewFileRepoWithBucket("kt-b")
	file, closer, err := repo.GetFileViaTemp(ctx, filepath.Join(id.String(), id.String()))
	if err != nil {
		return fmt.Errorf("get kt file: %w", err)
	}
	defer closer()
	if _, err := splitter.PrepareDicom(ctx, repo, id, file, splitter.Mp4); err != nil {
		return fmt.Errorf("prepare kt dicom: %w", err)
	}
	if err := s.dao.NewKtQuery(ctx).UpdateKtStatus(id, string(domain.ExamStatusPending)); err != nil {
		return fmt.Errorf("update kt status: %w", err)
	}
	if err := s.dbus.SendKtPrepared(ctx, &ktprepared.KtPrepared{
		KtId: id.String(),
	}); err != nil {
		return fmt.Errorf("send to ktprepared topic: %w", err)
	}

	return nil
}
