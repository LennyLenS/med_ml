package kt

import (
	"exam/internal/repository/kt/entity"
)

func (q *repo) InsertKt(kt entity.KT) error {
	query := q.QueryBuilder().
		Insert(table).
		Columns(
			columnID,
			columnChecked,
			columnAuthor,
			columnDeviceID,
			columnStatus,
			columnDescription,
			columnCreateAt,
		).
		Values(
			kt.Id,
			kt.Checked,
			kt.Author,
			kt.DeviceID,
			kt.Status,
			kt.Description,
			kt.CreateAt,
		)

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}

	return nil
}
