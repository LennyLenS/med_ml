package mappers

import (
	"github.com/google/uuid"

	domain "composition-api/internal/domain/exam"
	pb "composition-api/internal/generated/grpc/clients/exam"
)

type Image struct{}

func (m Image) Domain(pb *pb.Image) domain.Image {
	return domain.Image{
		Id:    uuid.MustParse(pb.Id),
		MriID: uuid.MustParse(pb.MriId),
		Page:  int(pb.Page),
	}
}

func (m Image) SliceDomain(pbs []*pb.Image) []domain.Image {
	return slice(pbs, m)
}


