package mappers

import (
	domain "composition-api/internal/domain/exam"
	api "composition-api/internal/generated/http/api"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestKtResponseContainsMLResult(t *testing.T) {
	input := domain.KT{Id: uuid.New(), Author: uuid.New(), CreateAt: time.Now(), Status: domain.ExamStatusCompleted, PredictedClasses: map[string]float64{"benign": 0.7, "malignant": 0.3}}
	output := (KT{}).Domain(input)
	require.NoError(t, output.Validate())
	require.Equal(t, api.KtStatusCompleted, output.Status)
	require.Equal(t, api.KtClassProbabilities(input.PredictedClasses), output.ClassProbabilities)
	input.Status = domain.ExamStatusNew
	input.PredictedClasses = nil
	output = (KT{}).Domain(input)
	require.NoError(t, output.Validate())
	require.NotNil(t, output.ClassProbabilities)
	require.Empty(t, output.ClassProbabilities)
}
