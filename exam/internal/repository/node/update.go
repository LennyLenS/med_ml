package node

import (
	"exam/internal/repository/node/entity"

	sq "github.com/Masterminds/squirrel"
)

func (q *repo) UpdateNode(node entity.Node) error {
	query := q.QueryBuilder().
		Update(table).
		SetMap(sq.Eq{
			columnValidation: node.Validation,
			columnKnosp012:   node.Knosp012,
			columnKnosp3:     node.Knosp3,
			columnKnosp4:     node.Knosp4,
		}).
		Where(sq.Eq{
			columnID: node.Id,
		})

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}

	return nil
}
