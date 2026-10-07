package mri

import (
	"exam/internal/repository/mri/entity"
)

func (q *repo) InsertMri(mri entity.Mri) error {
	query := q.QueryBuilder().
		Insert(table).
		Columns(
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
		Values(
			mri.Id,
			mri.Projection,
			mri.Checked,
			mri.ExternalID,
			mri.Author,
			mri.DeviceID,
			mri.Status,
			mri.Description,
			mri.CreateAt,
		)

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}

	return nil
}
