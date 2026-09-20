package mri

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"exam/internal/repository/mri/entity"
)

func (q *repo) UpdateMri(mri entity.Mri) error {
	query := q.QueryBuilder().
		Update(table).
		SetMap(sq.Eq{
			columnProjection: mri.Projection,
			columnChecked:    mri.Checked,
			columnStatus:     mri.Status,
		}).
		Where(sq.Eq{
			columnID: mri.Id,
		})

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}

	return nil
}

func (q *repo) UpdateMriStatus(id uuid.UUID, status string) error {
	query := q.QueryBuilder().
		Update(table).
		SetMap(sq.Eq{
			columnStatus: status,
		}).
		Where(sq.Eq{
			columnID: id,
		})

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}

	return nil
}
