package mri

import (
	"context"
	"fmt"
	"time"

	"exam/internal/domain"
	echographicEntity "exam/internal/repository/echographic/entity"
	mriEntity "exam/internal/repository/mri/entity"

	"github.com/google/uuid"
)

func (s *service) CreateMri(ctx context.Context, arg CreateMriArg) (uuid.UUID, error) {
	ctx, err := s.dao.BeginTx(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = s.dao.RollbackTx(ctx) }()

	mri := domain.Mri{
		Id:          uuid.New(),
		Projection:  arg.Projection,
		Checked:     false,
		ExternalID:  arg.ExternalID,
		Author:      arg.Author,
		DeviceID:    arg.DeviceID,
		Status:      domain.ExamStatusNew,
		Description: arg.Description,
		CreateAt:    time.Now(),
	}

	if err := s.dao.NewMriQuery(ctx).InsertMri(mriEntity.Mri{}.FromDomain(mri)); err != nil {
		return uuid.Nil, fmt.Errorf("insert mri: %w", err)
	}

	if err := s.dao.NewEchographicQuery(ctx).InsertEchographic(echographicEntity.Echographic{Id: mri.Id}); err != nil {
		return uuid.Nil, fmt.Errorf("insert echographic: %w", err)
	}

	if err := s.dao.CommitTx(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit transaction: %w", err)
	}

	return mri.Id, nil
}
