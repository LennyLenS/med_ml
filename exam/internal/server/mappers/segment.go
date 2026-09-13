package mappers

import (
	"exam/internal/domain"
	pb "exam/internal/generated/grpc/service"
)

func SegmentFromDomain(domain domain.Segment) *pb.Segment {
	return &pb.Segment{
		Id:        domain.Id.String(),
		ImageId:   domain.ImageID.String(),
		NodeId:    domain.NodeID.String(),
		Contor:    domain.Contor,
		Ai:        domain.Ai,
		Knosp_012: domain.Knosp012,
		Knosp_3:   domain.Knosp3,
		Knosp_4:   domain.Knosp4,
	}
}

func SliceSegmentFromDomain(domains []domain.Segment) []*pb.Segment {
	pbs := make([]*pb.Segment, 0, len(domains))
	for _, d := range domains {
		pbs = append(pbs, SegmentFromDomain(d))
	}
	return pbs
}
