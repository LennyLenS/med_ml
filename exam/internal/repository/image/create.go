package image

import (
	"exam/internal/repository/image/entity"
)

func (q *repo) InsertImages(images ...entity.Image) error {
	query := q.QueryBuilder().
		Insert(table).
		Columns(
			columnId,
			columnMriId,
			columnPage,
		)

	for _, v := range images {
		query = query.Values(
			v.Id,
			v.MriID,
			v.Page,
		)
	}

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}

	return nil
}
