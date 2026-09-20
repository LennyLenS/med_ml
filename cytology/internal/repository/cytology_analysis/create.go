package cytology_analysis

import (
	"cytology/internal/repository/cytology_analysis/entity"
	repoEntity "cytology/internal/repository/entity"
)

func (q *repo) Insert(analysis entity.CytologyAnalysis) error {
	query := q.QueryBuilder().
		Insert(table).
		Columns("id", "cytology_id", "original_image_id", "created_at").
		Values(analysis.ID, analysis.CytologyID, analysis.OriginalImageID, analysis.CreatedAt)

	_, err := q.Runner().Execx(q.Context(), query)
	return repoEntity.WrapDBError(err)
}
