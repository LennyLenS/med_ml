package entity

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClassProbabilitiesScan(t *testing.T) {
	for _, value := range []any{[]byte(`{"benign":0.7,"malignant":0.3}`), `{"benign":0.7,"malignant":0.3}`} {
		var probabilities ClassProbabilities
		require.NoError(t, probabilities.Scan(value))
		require.Equal(t, ClassProbabilities{"benign": 0.7, "malignant": 0.3}, probabilities)
		require.NoError(t, probabilities.Scan(nil))
		require.Nil(t, probabilities)
	}
	var probabilities ClassProbabilities
	require.Error(t, probabilities.Scan([]byte("broken")))
	require.Error(t, probabilities.Scan(42))
}
