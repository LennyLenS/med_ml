package exam

import (
	"context"

	"github.com/google/uuid"

	"composition-api/internal/adapters/exam/mappers"
	domain "composition-api/internal/domain/exam"
	pb "composition-api/internal/generated/grpc/clients/exam"
)

func (a *adapter) GetImagesByMriId(ctx context.Context, id uuid.UUID) ([]domain.Image, error) {
	res, err := a.client.GetImagesByMriId(ctx, &pb.GetImagesByMriIdIn{MriId: id.String()})
	if err != nil {
		return nil, err
	}

	return mappers.Image{}.SliceDomain(res.Images), nil
}


