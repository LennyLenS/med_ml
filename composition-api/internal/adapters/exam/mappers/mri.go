package mappers

import (
	"time"

	"github.com/google/uuid"

	domain "composition-api/internal/domain/exam"
	pb "composition-api/internal/generated/grpc/clients/exam"
)

var examStatusMap = map[pb.ExamStatus]domain.ExamStatus{
	pb.ExamStatus_EXAM_STATUS_NEW:       domain.ExamStatusNew,
	pb.ExamStatus_EXAM_STATUS_PENDING:   domain.ExamStatusPending,
	pb.ExamStatus_EXAM_STATUS_COMPLETED: domain.ExamStatusCompleted,
}

var mriProjectionMap = map[pb.MriProjection]domain.MriProjection{
	pb.MriProjection_MRI_PROJECTION_CROSS: domain.MriProjectionCross,
	pb.MriProjection_MRI_PROJECTION_LONG:  domain.MriProjectionLong,
}

type Mri struct{}

func (m Mri) Domain(pb *pb.Mri) domain.Mri {
	createAt, _ := time.Parse(time.RFC3339, pb.CreateAt)

	return domain.Mri{
		Id:          uuid.MustParse(pb.Id),
		Projection:  mriProjectionMap[pb.Projection],
		Checked:     pb.Checked,
		ExternalID:  uuid.MustParse(pb.ExternalId),
		Author:      uuid.MustParse(pb.Author),
		DeviceID:    int(pb.DeviceId),
		Status:      examStatusMap[pb.Status],
		Description: pb.Description,
		CreateAt:    createAt,
	}
}

func (m Mri) SliceDomain(pbs []*pb.Mri) []domain.Mri {
	return slice(pbs, m)
}


