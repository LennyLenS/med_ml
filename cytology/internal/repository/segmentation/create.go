package segmentation

import (
	repoEntity "cytology/internal/repository/entity"
	"cytology/internal/repository/segmentation/entity"
)

func (q *repo) InsertSegmentation(seg entity.Segmentation) error {
	// Вставляем сегментацию
	query := q.QueryBuilder().
		Insert(table).
		Columns(
			columnID,
			columnSegmentationGroupID,
			columnCreateAt,
		).
		Values(
			seg.Id,
			seg.SegmentationGroupID,
			seg.CreateAt,
		)

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return repoEntity.WrapDBError(err)
	}

	// Вставляем точки
	if len(seg.Points) > 0 {
		pointsQuery := q.QueryBuilder().
			Insert(pointTable).
			Columns(
				pointColumnID,
				pointColumnSegmentationID,
				pointColumnX,
				pointColumnY,
				pointColumnUID,
				pointColumnCreateAt,
			)

		for _, point := range seg.Points {
			pointsQuery = pointsQuery.Values(
				point.Id,
				point.SegmentationID,
				point.X,
				point.Y,
				point.UID,
				point.CreateAt,
			)
		}

		_, err = q.Runner().Execx(q.Context(), pointsQuery)
		if err != nil {
			return repoEntity.WrapDBError(err)
		}
	}

	return nil
}
