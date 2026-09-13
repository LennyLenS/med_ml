package node

import (
	"database/sql"
	"errors"
	"fmt"

	daoEntity "exam/internal/repository/entity"
	"exam/internal/repository/node/entity"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (q *repo) GetNodeByID(id uuid.UUID) (entity.Node, error) {
	query := q.QueryBuilder().
		Select(
			columnID,
			columnAI,
			columnMriID,
			columnValidation,
			columnKnosp012,
			columnKnosp3,
			columnKnosp4,
			columnDescription,
		).
		From(table).
		Where(sq.Eq{
			columnID: id,
		})

	var node entity.Node
	if err := q.Runner().Getx(q.Context(), &node, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Node{}, daoEntity.ErrNotFound
		}
		return entity.Node{}, err
	}

	return node, nil
}

func (q *repo) GetNodesByImageID(id uuid.UUID) ([]entity.Node, error) {
	query := q.QueryBuilder().
		Select(
			fmt.Sprintf("%s.%s", table, columnID),
			fmt.Sprintf("%s.%s", table, columnAI),
			fmt.Sprintf("%s.%s", table, columnMriID),
			fmt.Sprintf("%s.%s", table, columnValidation),
			fmt.Sprintf("%s.%s", table, columnKnosp012),
			fmt.Sprintf("%s.%s", table, columnKnosp3),
			fmt.Sprintf("%s.%s", table, columnKnosp4),
			fmt.Sprintf("%s.%s", table, columnDescription),
		).
		From(table). // TODO: вынести константы таблиц в отдельный пакет, тут пересечение с segment
		InnerJoin("segment ON segment.node_id = node.id").
		InnerJoin("image ON image.id = segment.image_id").
		Where(sq.Eq{
			"image.id": id,
		})

	var mri []entity.Node
	if err := q.Runner().Selectx(q.Context(), &mri, query); err != nil {
		return nil, err
	}

	if len(mri) == 0 {
		return nil, daoEntity.ErrNotFound
	}

	return mri, nil
}

func (q *repo) GetNodesByMriID(id uuid.UUID) ([]entity.Node, error) {
	query := q.QueryBuilder().
		Select(
			columnID,
			columnAI,
			columnMriID,
			columnValidation,
			columnKnosp012,
			columnKnosp3,
			columnKnosp4,
			columnDescription,
		).
		From(table).
		Where(sq.Eq{
			columnMriID: id,
		})

	var mri []entity.Node
	if err := q.Runner().Selectx(q.Context(), &mri, query); err != nil {
		return nil, err
	}

	if len(mri) == 0 {
		return nil, daoEntity.ErrNotFound
	}

	return mri, nil
}
