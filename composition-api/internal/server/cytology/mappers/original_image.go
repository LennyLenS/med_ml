package mappers

import (
	"net/url"

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
			// UUID не преобразуется в int, оставляем пустым
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
			// Создаем URL из ImagePath
			Set: img.ImagePath != "",
		},
	}

	if img.ImagePath != "" {
		// Формат: /download/cytology/{cytology_id}/{original_image_id}
		imageURLStr := "/download/cytology/" + img.CytologyID.String() + "/" + img.Id.String()
		imageURL, err := url.Parse(imageURLStr)
		if err == nil {
			result.Image = api.OptURI{
				Value: *imageURL,
				Set:   true,
			}
		}
	}

	if img.DelayTime != nil {
		result.DelayTime = api.OptFloat64{
			Value: *img.DelayTime,
			Set:   true,
		}
	}

	return result
}
