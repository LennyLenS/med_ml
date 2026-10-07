package image

import (
	daoEntity "exam/internal/repository/entity"
	"exam/internal/repository/image/entity"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (q *repo) GetImagesByMriID(mriID uuid.UUID) ([]entity.Image, error) {
	query := q.QueryBuilder().
		Select(
			columnId,
			columnMriId,
			columnPage,
		).
		From(table).
		Where(sq.Eq{
			columnMriId: mriID,
		})

	var images []entity.Image
	if err := q.Runner().Selectx(q.Context(), &images, query); err != nil {
		return nil, err
	}

	if len(images) == 0 {
		return nil, daoEntity.ErrNotFound
	}

	return images, nil
}
