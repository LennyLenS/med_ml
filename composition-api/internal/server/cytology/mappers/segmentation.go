package mappers

import (
	"github.com/google/uuid"

	domain "composition-api/internal/domain/cytology"
	api "composition-api/internal/generated/http/api"
	cytologySrv "composition-api/internal/services/cytology"
)

type SegmentationGroup struct{}

// Удалены неиспользуемые методы Domain, SliceDomain, CreateArg, UpdateArg - заменены на новые методы для работы с обновленными типами API

type Segmentation struct{}

// Удалены неиспользуемые методы Domain, SliceDomain, CreateArg, UpdateArg - заменены на новые методы для работы с обновленными типами API

func (SegmentationGroup) ToSegmentationDataList(groups []domain.SegmentationGroup) []api.CytologySegmentsListOKResultsItem {
	// TODO: Реализовать правильный маппинг с получением сегментов для каждой группы
	// Пока возвращаем упрощенную версию
	result := make([]api.CytologySegmentsListOKResultsItem, 0, len(groups))
	for _, group := range groups {
		item := api.CytologySegmentsListOKResultsItem{
			ID: api.OptInt{
				// TODO: Преобразовать UUID в int
				Set: false,
			},
			Data: []api.CytologySegmentsListOKResultsItemDataItem{},
			GroupType: api.OptCytologySegmentsListOKResultsItemGroupType{
				Value: api.CytologySegmentsListOKResultsItemGroupType(group.GroupType),
				Set:   true,
			},
			SegType: api.OptCytologySegmentsListOKResultsItemSegType{
				Value: api.CytologySegmentsListOKResultsItemSegType(group.SegType),
				Set:   true,
			},
			IsAi: api.OptBool{
				Value: group.IsAI,
				Set:   true,
			},
		}
		if group.Details != nil {
			item.Details = &api.CytologySegmentsListOKResultsItemDetails{}
		}
		result = append(result, item)
	}
	return result
}

func (SegmentationGroup) CreateArgFromCytologySegmentGroupCreateCreateReq(cytologyID uuid.UUID, req *api.CytologySegmentGroupCreateCreateReq) cytologySrv.CreateSegmentationGroupArg {
	// В swagger.json req содержит data (с points) и seg_type
	// Нужно определить group_type и is_ai из контекста или использовать значения по умолчанию
	arg := cytologySrv.CreateSegmentationGroupArg{
		CytologyID: cytologyID,
		SegType:    domain.SegType(req.SegType),
		GroupType:  domain.GroupType("CE"), // TODO: Определить из контекста или другого источника
		IsAI:       false,                  // TODO: Определить из контекста
	}

	return arg
}

func (SegmentationGroup) ToCytologySegmentGroupCreateCreateCreatedDataPoints(points []api.CytologySegmentGroupCreateCreateReqDataPointsItem) []api.CytologySegmentGroupCreateCreateCreatedDataPointsItem {
	result := make([]api.CytologySegmentGroupCreateCreateCreatedDataPointsItem, 0, len(points))
	for _, point := range points {
		result = append(result, api.CytologySegmentGroupCreateCreateCreatedDataPointsItem{
			X: point.X,
			Y: point.Y,
		})
	}
	return result
}

func (Segmentation) ToCytologySegmentUpdateReadOK(seg domain.Segmentation) api.CytologySegmentUpdateReadOK {
	result := api.CytologySegmentUpdateReadOK{
		Points: make([]api.CytologySegmentUpdateReadOKPointsItem, 0, len(seg.Points)),
		SegmentGroup: api.OptInt{
			// TODO: Преобразовать UUID в int
			Set: false,
		},
	}

	for _, point := range seg.Points {
		result.Points = append(result.Points, api.CytologySegmentUpdateReadOKPointsItem{
			X: point.X,
			Y: point.Y,
		})
	}

	return result
}

func (Segmentation) UpdateArgFromCytologySegmentUpdateUpdateReq(id uuid.UUID, req *api.CytologySegmentUpdateUpdateReq) cytologySrv.UpdateSegmentationArg {
	points := make([]domain.SegmentationPoint, 0, len(req.Points))
	for _, point := range req.Points {
		points = append(points, domain.SegmentationPoint{
			X: point.X,
			Y: point.Y,
		})
	}

	return cytologySrv.UpdateSegmentationArg{
		Id:     id,
		Points: points,
	}
}

func (Segmentation) UpdateArgFromCytologySegmentUpdatePartialUpdateReq(id uuid.UUID, req *api.CytologySegmentUpdatePartialUpdateReq) cytologySrv.UpdateSegmentationArg {
	points := make([]domain.SegmentationPoint, 0, len(req.Points))
	for _, point := range req.Points {
		points = append(points, domain.SegmentationPoint{
			X: point.X,
			Y: point.Y,
		})
	}

	return cytologySrv.UpdateSegmentationArg{
		Id:     id,
		Points: points,
	}
}

func (Segmentation) ToCytologySegmentUpdateUpdateOK(seg domain.Segmentation, req *api.CytologySegmentUpdateUpdateReq) api.CytologySegmentUpdateUpdateOK {
	result := api.CytologySegmentUpdateUpdateOK{
		Points: make([]api.CytologySegmentUpdateUpdateOKPointsItem, 0, len(seg.Points)),
		SegmentGroup: api.OptInt{
			// TODO: Преобразовать UUID в int
			Set: false,
		},
	}

	for _, point := range seg.Points {
		result.Points = append(result.Points, api.CytologySegmentUpdateUpdateOKPointsItem{
			X: point.X,
			Y: point.Y,
		})
	}

	return result
}

func (Segmentation) ToCytologySegmentUpdatePartialUpdateOK(seg domain.Segmentation, req *api.CytologySegmentUpdatePartialUpdateReq) api.CytologySegmentUpdatePartialUpdateOK {
	result := api.CytologySegmentUpdatePartialUpdateOK{
		Points: make([]api.CytologySegmentUpdatePartialUpdateOKPointsItem, 0, len(seg.Points)),
		SegmentGroup: api.OptInt{
			// TODO: Преобразовать UUID в int
			Set: false,
		},
	}

	for _, point := range seg.Points {
		result.Points = append(result.Points, api.CytologySegmentUpdatePartialUpdateOKPointsItem{
			X: point.X,
			Y: point.Y,
		})
	}

	return result
}
