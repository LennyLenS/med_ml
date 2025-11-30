package mappers

import (
	"github.com/google/uuid"

	api "composition-api/internal/generated/http/api"
	domain "composition-api/internal/domain/cytology"
	cytologySrv "composition-api/internal/services/cytology"
)

type SegmentationGroup struct{}

func (SegmentationGroup) Domain(group domain.SegmentationGroup) api.SegmentationGroup {
	result := api.SegmentationGroup{
		ID:         group.Id,
		CytologyID: group.CytologyID,
		SegType:    api.SegmentationGroupSegType(group.SegType),
		GroupType:  api.SegmentationGroupGroupType(group.GroupType),
		IsAi:       group.IsAI,
		CreateAt:   group.CreateAt,
	}

	if group.Details != nil {
		// Details is a JSON string, might need to parse it
		// For now, we'll set it to nil
		result.Details = nil
	}

	return result
}

func (SegmentationGroup) SliceDomain(groups []domain.SegmentationGroup) []api.SegmentationGroup {
	result := make([]api.SegmentationGroup, 0, len(groups))
	for _, group := range groups {
		result = append(result, SegmentationGroup{}.Domain(group))
	}
	return result
}

func (SegmentationGroup) CreateArg(cytologyID uuid.UUID, req *api.CytologyIDSegmentationGroupsPostReq) cytologySrv.CreateSegmentationGroupArg {
	isAI := false
	if req.IsAi.Set {
		isAI = req.IsAi.Value
	}

	arg := cytologySrv.CreateSegmentationGroupArg{
		CytologyID: cytologyID,
		SegType:    domain.SegType(req.SegType),
		GroupType:  domain.GroupType(req.GroupType),
		IsAI:       isAI,
	}

	if req.Details != nil {
		// Details handling - might need to marshal to JSON string
	}

	return arg
}

func (SegmentationGroup) UpdateArg(id uuid.UUID, req *api.CytologySegmentationGroupIDPatchReq) cytologySrv.UpdateSegmentationGroupArg {
	arg := cytologySrv.UpdateSegmentationGroupArg{
		Id: id,
	}

	if req.SegType.Set {
		segType := domain.SegType(req.SegType.Value)
		arg.SegType = &segType
	}

	if req.Details != nil {
		// Details handling
	}

	return arg
}

type Segmentation struct{}

func (Segmentation) Domain(seg domain.Segmentation) api.Segmentation {
	result := api.Segmentation{
		ID:                  seg.Id,
		SegmentationGroupID: seg.SegmentationGroupID,
		Points:              make([]api.SegmentationPoint, 0, len(seg.Points)),
		CreateAt:            seg.CreateAt,
	}

	for _, point := range seg.Points {
		result.Points = append(result.Points, api.SegmentationPoint{
			X: point.X,
			Y: point.Y,
		})
	}

	return result
}

func (Segmentation) SliceDomain(segs []domain.Segmentation) []api.Segmentation {
	result := make([]api.Segmentation, 0, len(segs))
	for _, seg := range segs {
		result = append(result, Segmentation{}.Domain(seg))
	}
	return result
}

func (Segmentation) CreateArg(groupID uuid.UUID, req *api.CytologySegmentationGroupIDSegmentsPostReq) cytologySrv.CreateSegmentationArg {
	points := make([]domain.SegmentationPoint, 0, len(req.Points))
	for _, point := range req.Points {
		points = append(points, domain.SegmentationPoint{
			X: point.X,
			Y: point.Y,
		})
	}

	return cytologySrv.CreateSegmentationArg{
		SegmentationGroupID: groupID,
		Points:              points,
	}
}

func (Segmentation) UpdateArg(id uuid.UUID, req *api.CytologySegmentationIDPatchReq) cytologySrv.UpdateSegmentationArg {
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
