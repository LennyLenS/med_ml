package image

import (
	"context"

	"github.com/AlekSi/pointer"

	api "composition-api/internal/generated/http/api"
	mappers "composition-api/internal/server/exam/mappers"
)

func (h *handler) MriIDImagesGet(ctx context.Context, params api.MriIDImagesGetParams) (api.MriIDImagesGetRes, error) {
	images, err := h.services.ExamImageService.GetImagesByMriID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return pointer.To(api.MriIDImagesGetOKApplicationJSON(mappers.SliceImage(images))), nil
}
