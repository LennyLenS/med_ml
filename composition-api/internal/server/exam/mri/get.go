package mri

import (
	"context"

	"github.com/AlekSi/pointer"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/exam/mappers"
)

func (h *handler) MriIDGet(ctx context.Context, params api.MriIDGetParams) (api.MriIDGetRes, error) {
	mri, err := h.services.MriService.GetByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return pointer.To(mappers.Mri{}.Domain(mri)), nil
}

func (h *handler) MrisExternalIDGet(ctx context.Context, params api.MrisExternalIDGetParams) (api.MrisExternalIDGetRes, error) {
	mris, err := h.services.MriService.GetByExternalID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return pointer.To(api.MrisExternalIDGetOKApplicationJSON(mappers.Mri{}.SliceDomain(mris))), nil
}

func (h *handler) MrisAuthorIDGet(ctx context.Context, params api.MrisAuthorIDGetParams) (api.MrisAuthorIDGetRes, error) {
	mris, err := h.services.MriService.GetByAuthor(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return pointer.To(api.MrisAuthorIDGetOKApplicationJSON(mappers.Mri{}.SliceDomain(mris))), nil
}

func (h *handler) MriIDEchographicsGet(ctx context.Context, params api.MriIDEchographicsGetParams) (api.MriIDEchographicsGetRes, error) {
	echographics, err := h.services.MriService.GetEchographicsByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return pointer.To(mappers.Echographics(echographics)), nil
}
