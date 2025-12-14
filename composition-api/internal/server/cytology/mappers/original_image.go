package mappers

import (
	domain "composition-api/internal/domain/cytology"
	api "composition-api/internal/generated/http/api"
)

type OriginalImage struct{}

// Методы Domain и SliceDomain удалены, так как тип api.OriginalImage не существует в сгенерированном API
// Используются только методы для работы с CytologyReadOKOriginalImage

func (OriginalImage) ToCytologyReadOKOriginalImage(img *domain.OriginalImage) api.CytologyReadOKOriginalImage {
	if img == nil {
		return api.CytologyReadOKOriginalImage{}
	}

	result := api.CytologyReadOKOriginalImage{
		ID: api.OptInt{
			// TODO: Преобразовать UUID в int
			Set: false,
		},
		CreateDate: api.OptDateTime{
			Value: img.CreateDate,
			Set:   true,
		},
		ViewedFlag: api.OptBool{
			Value: img.ViewedFlag,
			Set:   true,
		},
		Image: api.OptURI{
			// TODO: Создать URL из ImagePath
			Set: false,
		},
	}

	if img.DelayTime != nil {
		result.DelayTime = api.OptFloat64{
			Value: *img.DelayTime,
			Set:   true,
		}
	}

	return result
}
