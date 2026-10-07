package mappers

import (
	domain "composition-api/internal/domain/exam"
	api "composition-api/internal/generated/http/api"
)

type Segment struct{}

func (Segment) Domain(segment domain.Segment) (api.ExamSegment, error) {
	contor, err := Contor(segment.Contor)
	if err != nil {
		return api.ExamSegment{}, err
	}

	return api.ExamSegment{
		ID:       segment.Id,
		ImageID:  segment.ImageID,
		NodeID:   segment.NodeID,
		Contor:   contor,
		Ai:       segment.Ai,
		Knosp012: segment.Knosp012,
		Knosp3:   segment.Knosp3,
		Knosp4:   segment.Knosp4,
	}, nil
}

func (Segment) SliceDomain(segments []domain.Segment) ([]api.ExamSegment, error) {
	result := make([]api.ExamSegment, 0, len(segments))
	for _, segment := range segments {
		segment, err := Segment{}.Domain(segment)
		if err != nil {
			return nil, err
		}
		result = append(result, segment)
	}
	return result, nil
}
