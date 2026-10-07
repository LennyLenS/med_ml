package mappers

import (
	domain "composition-api/internal/domain/exam"
	api "composition-api/internal/generated/http/api"
)

func Image(image domain.Image) api.ExamImage {
	return api.ExamImage{
		ID:    image.Id,
		MriID: image.MriID,
		Page:  image.Page,
	}
}

func SliceImage(images []domain.Image) []api.ExamImage {
	result := make([]api.ExamImage, 0, len(images))
	for _, image := range images {
		result = append(result, Image(image))
	}
	return result
}
