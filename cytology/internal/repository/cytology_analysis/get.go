package cytology_analysis

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"cytology/internal/repository/cytology_analysis/entity"
	repoEntity "cytology/internal/repository/entity"
)

func (q *repo) GetByID(id uuid.UUID) (entity.CytologyAnalysis, error) {
	query := q.QueryBuilder().
		Select("id", "cytology_id", "original_image_id", "result", "created_at", "completed_at").
		From(table).
		Where(sq.Eq{"id": id})

	var analysis entity.CytologyAnalysis
	if err := q.Runner().Getx(q.Context(), &analysis, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.CytologyAnalysis{}, repoEntity.ErrNotFound
		}
		return entity.CytologyAnalysis{}, err
	}
	return analysis, nil
}
