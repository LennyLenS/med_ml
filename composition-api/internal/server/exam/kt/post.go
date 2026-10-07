package kt

import (
	"context"
	"fmt"

	"composition-api/internal/generated/http/api"
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

	contentType := req.File.Header.Get("Content-Type")
	if contentType != "video/mp4" {
		return nil, fmt.Errorf("wrong file type, expected: mp4, got: %s", contentType)
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
