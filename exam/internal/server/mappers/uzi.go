package mappers

import (
	"time"

	"exam/internal/domain"
	pb "exam/internal/generated/grpc/service"
)

var examStatusMap = map[domain.ExamStatus]pb.ExamStatus{
	domain.ExamStatusNew:       pb.ExamStatus_EXAM_STATUS_NEW,
	domain.ExamStatusPending:   pb.ExamStatus_EXAM_STATUS_PENDING,
	domain.ExamStatusCompleted: pb.ExamStatus_EXAM_STATUS_COMPLETED,
}

var MriProjectionMap = map[domain.MriProjection]pb.MriProjection{
	domain.MriProjectionLong:  pb.MriProjection_MRI_PROJECTION_LONG,
	domain.MriProjectionCross: pb.MriProjection_MRI_PROJECTION_CROSS,
}

var MriProjectionReverseMap = map[pb.MriProjection]domain.MriProjection{
	pb.MriProjection_MRI_PROJECTION_LONG:  domain.MriProjectionLong,
	pb.MriProjection_MRI_PROJECTION_CROSS: domain.MriProjectionCross,
}

func MriFromDomain(domain domain.Mri) *pb.Mri {
	return &pb.Mri{
		Id:          domain.Id.String(),
		Projection:  MriProjectionMap[domain.Projection],
		Checked:     domain.Checked,
		ExternalId:  domain.ExternalID.String(),
		Author:      domain.Author.String(),
		DeviceId:    int64(domain.DeviceID),
		Status:      examStatusMap[domain.Status],
		Description: domain.Description,
		CreateAt:    domain.CreateAt.Format(time.RFC3339),
	}
}

func SliceMriFromDomain(domains []domain.Mri) []*pb.Mri {
	pbs := make([]*pb.Mri, 0, len(domains))
	for _, d := range domains {
		pbs = append(pbs, MriFromDomain(d))
	}
	return pbs
}

func mapToProbabilities(classes map[string]float64) *pb.ClassProbabilities {
	return &pb.ClassProbabilities{
		ClassProbabilities: classes,
	}
}

func KtFromDomain(domain domain.KT) *pb.Kt {
	return &pb.Kt{
		Id:                 domain.Id.String(),
		Checked:            domain.Checked,
		Author:             domain.Author.String(),
		DeviceId:           int64(domain.DeviceID),
		Status:             examStatusMap[domain.Status],
		Description:        domain.Description,
		CreateAt:           domain.CreateAt.Format(time.RFC3339),
		ClassProbabilities: mapToProbabilities(domain.PredictedClasses),
	}
}

func SliceKtFromDomain(domains []domain.KT) []*pb.Kt {
	pbs := make([]*pb.Kt, 0, len(domains))
	for _, d := range domains {
		pbs = append(pbs, KtFromDomain(d))
	}
	return pbs
}
