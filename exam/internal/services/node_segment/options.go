package node_segment

import "exam/internal/domain"

// -- CreateNodesWithSegmentsOptions --

type createNodesWithSegmentsOption struct {
	newMriStatus       *domain.ExamStatus
	setNodesValidation *domain.NodeValidation
}

type CreateNodesWithSegmentsOption func(*createNodesWithSegmentsOption)

var (
	WithNewMriStatus = func(status domain.ExamStatus) CreateNodesWithSegmentsOption {
		return func(o *createNodesWithSegmentsOption) { o.newMriStatus = &status }
	}

	WithSetNodesValidation = func(validation domain.NodeValidation) CreateNodesWithSegmentsOption {
		return func(o *createNodesWithSegmentsOption) { o.setNodesValidation = &validation }
	}
)
