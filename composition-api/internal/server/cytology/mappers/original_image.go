package mappers

import (
	api "composition-api/internal/generated/http/api"
	domain "composition-api/internal/domain/cytology"
)

type OriginalImage struct{}

func (OriginalImage) Domain(img domain.OriginalImage) api.OriginalImage {
	result := api.OriginalImage{
		ID:         img.Id,
		CytologyID: img.CytologyID,
		ImagePath:  img.ImagePath,
		ViewedFlag: img.ViewedFlag,
	}

	if img.DelayTime != nil {
		result.DelayTime = api.OptFloat64{
			Value: *img.DelayTime,
			Set:   true,
		}
	}

	return result
}

func (OriginalImage) SliceDomain(imgs []domain.OriginalImage) []api.OriginalImage {
	result := make([]api.OriginalImage, 0, len(imgs))
	for _, img := range imgs {
		result = append(result, OriginalImage{}.Domain(img))
	}
	return result
}
