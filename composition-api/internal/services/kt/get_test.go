package kt

import (
	"context"
	"errors"
	"testing"

	"composition-api/internal/adapters"
	"composition-api/internal/adapters/exam"
	domain "composition-api/internal/domain/exam"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type resultAdapter struct {
	exam.Adapter
	id     uuid.UUID
	result domain.KT
	err    error
}

func (a *resultAdapter) GetKtById(ctx context.Context, id uuid.UUID) (domain.KT, error) {
	a.id = id
	return a.result, a.err
}

func TestGetByIDReturnsStoredResult(t *testing.T) {
	id := uuid.New()
	a := &resultAdapter{result: domain.KT{Id: id, Status: domain.ExamStatusCompleted, PredictedClasses: map[string]float64{"normal": 1}}}
	s := &service{adapters: &adapters.Adapters{Exam: a}}
	result, err := s.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, a.id)
	require.Equal(t, a.result, result)
	a.err = errors.New("unavailable")
	_, err = s.GetByID(context.Background(), id)
	require.ErrorIs(t, err, a.err)
}
