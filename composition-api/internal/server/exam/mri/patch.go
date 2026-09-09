package mri

import (
	"context"

	"github.com/AlekSi/pointer"

	domain "composition-api/internal/domain/exam"
	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/exam/mappers"
	apimappers "composition-api/internal/server/mappers"
	mriSrv "composition-api/internal/services/mri"
)

func (h *handler) MriIDPatch(ctx context.Context, req *api.MriIDPatchReq, params api.MriIDPatchParams) (api.MriIDPatchRes, error) {
	var projection *domain.MriProjection
	if req.Projection.IsSet() {
		projection = (*domain.MriProjection)(&req.Projection.Value)
	}

	mri, err := h.services.MriService.Update(ctx, mriSrv.UpdateMriArg{
		Id:         params.ID,
		Projection: projection,
		Checked:    apimappers.FromOptBool(req.Checked),
	})
	if err != nil {
		return nil, err
	}
	return pointer.To(mappers.Mri{}.Domain(mri)), nil
}

func (h *handler) MriIDEchographicsPatch(ctx context.Context, req *api.ExamEchographics, params api.MriIDEchographicsPatchParams) (api.MriIDEchographicsPatchRes, error) {
	echographics, err := h.services.MriService.UpdateEchographics(ctx, domain.Echographic{
		Id:              params.ID,
		Contors:         apimappers.FromOptString(req.Contors),
		LeftLobeLength:  apimappers.FromOptFloat64(req.LeftLobeLength),
		LeftLobeWidth:   apimappers.FromOptFloat64(req.LeftLobeWidth),
		LeftLobeThick:   apimappers.FromOptFloat64(req.LeftLobeThick),
		LeftLobeVolum:   apimappers.FromOptFloat64(req.LeftLobeVolum),
		RightLobeLength: apimappers.FromOptFloat64(req.RightLobeLength),
		RightLobeWidth:  apimappers.FromOptFloat64(req.RightLobeWidth),
		RightLobeThick:  apimappers.FromOptFloat64(req.RightLobeThick),
		RightLobeVolum:  apimappers.FromOptFloat64(req.RightLobeVolum),
		GlandVolum:      apimappers.FromOptFloat64(req.GlandVolum),
		Isthmus:         apimappers.FromOptFloat64(req.Isthmus),
		Struct:          apimappers.FromOptString(req.Struct),
		Echogenicity:    apimappers.FromOptString(req.Echogenicity),
		RegionalLymph:   apimappers.FromOptString(req.RegionalLymph),
		Vascularization: apimappers.FromOptString(req.Vascularization),
		Location:        apimappers.FromOptString(req.Location),
		Additional:      apimappers.FromOptString(req.Additional),
		Conclusion:      apimappers.FromOptString(req.Conclusion),
	})
	if err != nil {
		return nil, err
	}
	return pointer.To(mappers.Echographics(echographics)), nil
}
