package segmentation

import (
	"context"
	"errors"
	"net/http"

	"composition-api/internal/domain"
	domainCytology "composition-api/internal/domain/cytology"

	"github.com/AlekSi/pointer"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/cytology/mappers"
)

func (h *handler) CytologyIDSegmentationGroupsGet(ctx context.Context, params api.CytologyIDSegmentationGroupsGetParams) (api.CytologyIDSegmentationGroupsGetRes, error) {
	var segType *domainCytology.SegType
	var groupType *domainCytology.GroupType
	var isAI *bool

	if params.SegType.Set {
		t := domainCytology.SegType(params.SegType.Value)
		segType = &t
	}

	if params.GroupType.Set {
		t := domainCytology.GroupType(params.GroupType.Value)
		groupType = &t
	}

	if params.IsAi.Set {
		isAI = &params.IsAi.Value
	}

	groups, err := h.services.CytologyService.GetSegmentationGroupsByCytologyId(ctx, params.ID, segType, groupType, isAI)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &api.CytologyIDSegmentationGroupsGetNotFound{
				StatusCode: http.StatusNotFound,
				Response: api.Error{
					Message: "Группы сегментаций не найдены",
				},
			}, nil
		}
		return nil, err
	}

	return pointer.To(api.CytologyIDSegmentationGroupsGetOKApplicationJSON(
		mappers.SegmentationGroup{}.SliceDomain(groups),
	)), nil
}
