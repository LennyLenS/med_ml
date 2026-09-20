package segment

import (
	"exam/internal/repository/segment/entity"
)

func (q *repo) InsertSegments(segments ...entity.Segment) error {
	query := q.QueryBuilder().
		Insert(table).
		Columns(
			columnID,
			columnNodeID,
			columnImageID,
			columnContor,
			columnAi,
			columnKnosp012,
			columnKnosp3,
			columnKnosp4,
		)

	for _, v := range segments {
		query = query.Values(
			v.Id,
			v.NodeID,
			v.ImageID,
			v.Contor,
			v.Ai,
			v.Knosp012,
			v.Knosp3,
			v.Knosp4,
		)
	}
	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}

	return nil
}
