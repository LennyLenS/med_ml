package kt

import (
	"context"
	"fmt"

	"composition-api/internal/generated/http/api"
	"composition-api/internal/server/exam/upload"
	"composition-api/internal/server/mappers"
	"composition-api/internal/server/security"
	"composition-api/internal/services/kt"

	"github.com/AlekSi/pointer"
)

func (h *handler) KtPost(ctx context.Context, req *api.KtPostReq) (api.KtPostRes, error) {
	token, err := security.ParseToken(ctx)
	if err != nil {
		return nil, err
	}

	contentType := upload.NormalizeFormat(&req.File)
	if contentType != "video/mp4" && contentType != "application/dicom" {
		return nil, fmt.Errorf("wrong file type, expected: video/mp4 or application/dicom, got: %s", contentType)
	}

	ktID, err := h.services.KTService.Create(ctx, kt.CreateKtArg{
		File:        req.File,
		Author:      token.Id,
		DeviceID:    req.DeviceID,
		Description: mappers.FromOptString(req.Description),
	})
	if err != nil {
		return nil, err
	}

	return pointer.To(api.SimpleUuid{ID: ktID}), nil
}
