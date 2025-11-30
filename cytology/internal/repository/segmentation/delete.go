package segmentation

import (
	"github.com/google/uuid"

	sq "github.com/Masterminds/squirrel"

	repoEntity "cytology/internal/repository/entity"
)

func (r *repo) DeleteSegmentation(id uuid.UUID) error {
	// Точки удалятся каскадно из-за ON DELETE CASCADE
	query := r.QueryBuilder().
		Delete(table).
		Where(sq.Eq{
			columnID: id,
		})

	_, err := r.Runner().Execx(r.Context(), query)
	if err != nil {
		return repoEntity.WrapDBError(err)
	}

	return nil
}
