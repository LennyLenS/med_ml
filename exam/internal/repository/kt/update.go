package kt

import (
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"exam/internal/repository/kt/entity"
)

func (q *repo) UpdateKt(kt entity.KT) error {
	probabilities, err := serializeClassProbabilities(kt.PredictedClasses)
	if err != nil {
		return err
	}
	query := q.QueryBuilder().
		Update(table).
		SetMap(sq.Eq{
			columnChecked:          kt.Checked,
			columnStatus:           kt.Status,
			columnPredictedClasses: probabilities,
		}).
		Where(sq.Eq{
			columnID: kt.Id,
		})

	_, err = q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}

	return nil
}

func (q *repo) UpdateKtPrediction(id uuid.UUID, predictedClasses []byte) error {
	query := q.QueryBuilder().
		Update(table).
		SetMap(sq.Eq{
			columnPredictedClasses: predictedClasses,
		}).
		Where(sq.Eq{
			columnID: id,
		})

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}

	return nil
}

func serializeClassProbabilities(classProbabilities map[string]float64) ([]byte, error) {
	probsJSON, err := json.Marshal(classProbabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize: %v", err)
	}
	return probsJSON, err
}
