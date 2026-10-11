package mri

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"

	domain "composition-api/internal/domain/exam"
	api "composition-api/internal/generated/http/api"
	"composition-api/internal/server/exam/upload"
	mappers "composition-api/internal/server/mappers"
	"composition-api/internal/server/security"
	mriSrv "composition-api/internal/services/mri"
)

var mriProjectionMap = map[api.MriPostReqProjection]domain.MriProjection{
	api.MriPostReqProjectionCross: domain.MriProjectionCross,
	api.MriPostReqProjectionLong:  domain.MriProjectionLong,
}

func (h *handler) MriPost(ctx context.Context, req *api.MriPostReq) (api.MriPostRes, error) {
	token, err := security.ParseToken(ctx)
	if err != nil {
		return nil, err
	}

	contentType := upload.NormalizeFormat(&req.File)
	if contentType != "image/tiff" && contentType != "application/dicom" {
		return nil, fmt.Errorf("wrong file type, expected: image/tiff or application/dicom, got: %s", contentType)
	}

	mriID, err := h.services.MriService.Create(ctx, mriSrv.CreateMriArg{
		File:        req.File,
		Projection:  mriProjectionMap[req.Projection],
		ExternalID:  req.ExternalID,
		Author:      token.Id,
		DeviceID:    req.DeviceID,
		Description: mappers.FromOptString(req.Description),
	})
	if err != nil {
		return nil, err
	}

	return pointer.To(api.SimpleUuid{ID: mriID}), nil
}
