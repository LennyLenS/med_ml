package segment

import (
	"exam/internal/repository/segment/entity"

	sq "github.com/Masterminds/squirrel"
)

func (q *repo) UpdateSegment(segment entity.Segment) error {
	query := q.QueryBuilder().
		Update(table).
		SetMap(sq.Eq{
			columnContor:   segment.Contor,
			columnKnosp012: segment.Knosp012,
			columnKnosp3:   segment.Knosp3,
			columnKnosp4:   segment.Knosp4,
		}).
		Where(sq.Eq{
			columnID: segment.Id,
		})

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}

	return nil
}
