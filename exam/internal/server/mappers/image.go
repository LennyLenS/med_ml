package mappers

import (
	"exam/internal/domain"
	pb "exam/internal/generated/grpc/service"
)

func ImageFromDomain(domain domain.Image) *pb.Image {
	return &pb.Image{
		Id:    domain.Id.String(),
		MriId: domain.MriID.String(),
		Page:  int64(domain.Page),
	}
}

func SliceImageFromDomain(domains []domain.Image) []*pb.Image {
	pbs := make([]*pb.Image, 0, len(domains))
	for _, d := range domains {
		pbs = append(pbs, ImageFromDomain(d))
	}
	return pbs
}
