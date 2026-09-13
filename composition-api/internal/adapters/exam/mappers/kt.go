package mappers

import (
	"time"

	"github.com/google/uuid"

	domain "composition-api/internal/domain/exam"
	pb "composition-api/internal/generated/grpc/clients/exam"
)

type Kt struct{}

func (m Kt) Domain(pb *pb.Kt) domain.KT {
	createAt, _ := time.Parse(time.RFC3339, pb.CreateAt)

	return domain.KT{
		Id:               uuid.MustParse(pb.Id),
		Checked:          pb.Checked,
		Author:           uuid.MustParse(pb.Author),
		DeviceID:         int(pb.DeviceId),
		Status:           examStatusMap[pb.Status],
		Description:      pb.Description,
		CreateAt:         createAt,
		PredictedClasses: pb.ClassProbabilities.ClassProbabilities,
	}
}

func (m Kt) SliceDomain(pbs []*pb.Kt) []domain.KT {
	domains := make([]domain.KT, 0, len(pbs))
	for _, pb := range pbs {
		domains = append(domains, m.Domain(pb))
	}
	return domains
}

func MapToProbabilities(classes map[string]float64) *pb.ClassProbabilities {
	return &pb.ClassProbabilities{
		ClassProbabilities: classes,
	}
}


