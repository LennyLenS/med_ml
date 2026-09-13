package kt

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	daoEntity "exam/internal/repository/entity"
	"exam/internal/repository/kt/entity"
)

func (q *repo) GetKtByID(id uuid.UUID) (entity.KT, error) {
	query := q.QueryBuilder().
		Select(
			columnID,
			columnChecked,
			columnAuthor,
			columnDeviceID,
			columnStatus,
			columnDescription,
			columnCreateAt,
			columnPredictedClasses,
		).
		From(table).
		Where(sq.Eq{
			columnID: id,
		})

	var kt entity.KT
	if err := q.Runner().Getx(q.Context(), &kt, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.KT{}, daoEntity.ErrNotFound
		}
		return entity.KT{}, err
	}

	return kt, nil
}

func (q *repo) GetKtsByAuthor(author uuid.UUID) ([]entity.KT, error) {
	query := q.QueryBuilder().
		Select(
			columnID,
			columnChecked,
			columnAuthor,
			columnDeviceID,
			columnStatus,
			columnDescription,
			columnCreateAt,
			columnPredictedClasses,
		).
		From(table).
		Where(sq.Eq{
			columnAuthor: author,
		})

	var kt []entity.KT
	if err := q.Runner().Selectx(q.Context(), &kt, query); err != nil {
		return nil, err
	}

	if len(kt) == 0 {
		return nil, daoEntity.ErrNotFound
	}

	return kt, nil
}

func (q *repo) CheckExist(id uuid.UUID) (bool, error) {
	query := q.QueryBuilder().
		Select(
			columnID,
			columnChecked,
			columnAuthor,
			columnDeviceID,
			columnStatus,
			columnDescription,
			columnCreateAt,
			columnPredictedClasses,
		).
		Prefix("SELECT EXISTS (").
		From(table).
		Where(sq.Eq{
			columnID: id,
		}).
		Suffix(")")

	var exists bool
	if err := q.Runner().Getx(q.Context(), &exists, query); err != nil {
		return false, err
	}

	return exists, nil
}
