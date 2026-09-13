package node

import (
	"exam/internal/repository/node/entity"
)

func (q *repo) InsertNodes(nodes ...entity.Node) error {
	query := q.QueryBuilder().
		Insert(table).
		Columns(
			columnID,
			columnAI,
			columnMriID,
			columnValidation,
			columnKnosp012,
			columnKnosp3,
			columnKnosp4,
			columnDescription,
		)

	for _, v := range nodes {
		query = query.Values(
			v.Id,
			v.Ai,
			v.MriID,
			v.Validation,
			v.Knosp012,
			v.Knosp3,
			v.Knosp4,
			v.Description,
		)
	}
	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}

	return nil
}
