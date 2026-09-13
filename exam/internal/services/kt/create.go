package kt

import (
	"context"
	"fmt"
	"time"

	"exam/internal/domain"
	ktEntity "exam/internal/repository/kt/entity"

	"github.com/google/uuid"
)

func (s *service) CreateKt(ctx context.Context, arg CreateKTArg) (uuid.UUID, error) {
	ctx, err := s.dao.BeginTx(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = s.dao.RollbackTx(ctx) }()

	kt := domain.KT{
		Id:          uuid.New(),
		Checked:     false,
		Author:      arg.Author,
		DeviceID:    arg.DeviceID,
		Status:      domain.ExamStatusNew,
		Description: arg.Description,
		CreateAt:    time.Now(),
	}

	if err := s.dao.NewKtQuery(ctx).InsertKt(ktEntity.KT{}.FromDomain(kt)); err != nil {
		return uuid.Nil, fmt.Errorf("insert kt: %w", err)
	}

	if err := s.dao.CommitTx(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit transaction: %w", err)
	}

	return kt.Id, nil
}

func (s *service) SaveProcessedKt(ctx context.Context, ktID uuid.UUID, arg []byte) error {
	ktQuery := s.dao.NewKtQuery(ctx)
	if err := ktQuery.UpdateKtPrediction(ktID, arg); err != nil {
		return fmt.Errorf("update kt prediction: %w", err)
	}
	return nil
}
