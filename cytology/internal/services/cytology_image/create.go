package cytology_image

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cytology/internal/domain"
	cytologyImageEntity "cytology/internal/repository/cytology_image/entity"
	"cytology/internal/repository/entity"

	"github.com/google/uuid"
)

func (s *service) CreateCytologyImage(ctx context.Context, arg CreateCytologyImageArg) (uuid.UUID, error) {
	ctx, err := s.dao.BeginTx(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = s.dao.RollbackTx(ctx) }()

	img := domain.CytologyImage{
		Id:                uuid.New(),
		ExternalID:        arg.ExternalID,
		PatientCardID:     arg.PatientCardID,
		DiagnosticNumber:  arg.DiagnosticNumber,
		DiagnosticMarking: arg.DiagnosticMarking,
		MaterialType:      arg.MaterialType,
		DiagnosDate:       time.Now(),
		IsLast:            true,
		Calcitonin:        arg.Calcitonin,
		CalcitoninInFlush: arg.CalcitoninInFlush,
		Thyroglobulin:     arg.Thyroglobulin,
		Details:           arg.Details,
		PrevID:            arg.PrevID,
		ParentPrevID:      arg.ParentPrevID,
		CreateAt:          time.Now(),
	}

	if err := s.dao.NewCytologyImageQuery(ctx).InsertCytologyImage(cytologyImageEntity.CytologyImage{}.FromDomain(img)); err != nil {
		var valErr *entity.DBValidationError
		if errors.As(err, &valErr) {
			return uuid.Nil, domain.ErrUnprocessableEntity
		}
		return uuid.Nil, fmt.Errorf("insert cytology image: %w", err)
	}

	if err := s.dao.CommitTx(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit transaction: %w", err)
	}

	return img.Id, nil
}
