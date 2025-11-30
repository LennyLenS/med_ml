package segmentation_group

import (
	repoEntity "cytology/internal/repository/entity"
	"cytology/internal/repository/segmentation_group/entity"
)

func (q *repo) InsertSegmentationGroup(group entity.SegmentationGroup) error {
	query := q.QueryBuilder().
		Insert(table).
		Columns(
			columnID,
			columnCytologyID,
			columnSegType,
			columnGroupType,
			columnIsAI,
			columnDetails,
			columnCreateAt,
		).
		Values(
			group.Id,
			group.CytologyID,
			group.SegType,
			group.GroupType,
			group.IsAI,
			group.Details,
			group.CreateAt,
		)

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return repoEntity.WrapDBError(err)
	}

	return nil
}
