package mri

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	daoEntity "exam/internal/repository/entity"
	"exam/internal/repository/mri/entity"
)

func (q *repo) GetMriByID(id uuid.UUID) (entity.Mri, error) {
	query := q.QueryBuilder().
		Select(
			columnID,
			columnProjection,
			columnChecked,
			columnExternalID,
			columnAuthor,
			columnDeviceID,
			columnStatus,
			columnDescription,
			columnCreateAt,
		).
		From(table).
		Where(sq.Eq{
			columnID: id,
		})

	var mri entity.Mri
	if err := q.Runner().Getx(q.Context(), &mri, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Mri{}, daoEntity.ErrNotFound
		}
		return entity.Mri{}, err
	}

	return mri, nil
}

func (q *repo) GetMrisByExternalID(externalID uuid.UUID) ([]entity.Mri, error) {
	query := q.QueryBuilder().
		Select(
			columnID,
			columnProjection,
			columnChecked,
			columnExternalID,
			columnAuthor,
			columnDeviceID,
			columnStatus,
			columnDescription,
			columnCreateAt,
		).
		From(table).
		Where(sq.Eq{
			columnExternalID: externalID,
		})

	var mri []entity.Mri
	if err := q.Runner().Selectx(q.Context(), &mri, query); err != nil {
		return nil, err
	}

	if len(mri) == 0 {
		return nil, daoEntity.ErrNotFound
	}

	return mri, nil
}

func (q *repo) GetMrisByAuthor(author uuid.UUID) ([]entity.Mri, error) {
	query := q.QueryBuilder().
		Select(
			columnID,
			columnProjection,
			columnChecked,
			columnExternalID,
			columnAuthor,
			columnDeviceID,
			columnStatus,
			columnDescription,
			columnCreateAt,
		).
		From(table).
		Where(sq.Eq{
			columnAuthor: author,
		})

	var mri []entity.Mri
	if err := q.Runner().Selectx(q.Context(), &mri, query); err != nil {
		return nil, err
	}

	if len(mri) == 0 {
		return nil, daoEntity.ErrNotFound
	}

	return mri, nil
}

func (q *repo) CheckExist(id uuid.UUID) (bool, error) {
	query := q.QueryBuilder().
		Select(
			columnID,
			columnProjection,
			columnChecked,
			columnExternalID,
			columnAuthor,
			columnDeviceID,
			columnStatus,
			columnDescription,
			columnCreateAt,
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
