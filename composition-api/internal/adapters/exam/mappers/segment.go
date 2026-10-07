package mappers

import (
	"github.com/google/uuid"

	domain "composition-api/internal/domain/exam"
	pb "composition-api/internal/generated/grpc/clients/exam"
)

type Segment struct{}

func (m Segment) Domain(pb *pb.Segment) domain.Segment {
	return domain.Segment{
		Id:       uuid.MustParse(pb.Id),
		ImageID:  uuid.MustParse(pb.ImageId),
		NodeID:   uuid.MustParse(pb.NodeId),
		Contor:   pb.Contor,
		Ai:       pb.Ai,
		Knosp012: pb.Knosp_012,
		Knosp3:   pb.Knosp_3,
		Knosp4:   pb.Knosp_4,
	}
}

func (m Segment) SliceDomain(pbs []*pb.Segment) []domain.Segment {
	return slice(pbs, m)
}


