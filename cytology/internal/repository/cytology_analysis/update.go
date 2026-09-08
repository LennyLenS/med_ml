package cytology_analysis

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	repoEntity "cytology/internal/repository/entity"
)

func (q *repo) Complete(id uuid.UUID, result []byte) error {
	query := q.QueryBuilder().
		Update(table).
		Set("result", string(result)).
		Set("completed_at", time.Now()).
		Where(sq.Eq{"id": id})

	_, err := q.Runner().Execx(q.Context(), query)
	return repoEntity.WrapDBError(err)
}
