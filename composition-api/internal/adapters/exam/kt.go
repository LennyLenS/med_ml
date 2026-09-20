package exam

import (
	"context"

	"github.com/google/uuid"

	"composition-api/internal/adapters/exam/mappers"
	domain "composition-api/internal/domain/exam"
	pb "composition-api/internal/generated/grpc/clients/exam"
)

func (a *adapter) CreateKt(ctx context.Context, in CreateKtIn) (uuid.UUID, error) {
	res, err := a.client.CreateKt(ctx, &pb.CreateKtIn{
		Author:      in.Author.String(),
		DeviceId:    int64(in.DeviceID),
		Description: in.Description,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.MustParse(res.Id), nil
}

func (a *adapter) GetKtById(ctx context.Context, id uuid.UUID) (domain.KT, error) {
	res, err := a.client.GetKtById(ctx, &pb.GetKtByIdIn{Id: id.String()})
	if err != nil {
		return domain.KT{}, err
	}

	return mappers.Kt{}.Domain(res.Kt), nil
}

func (a *adapter) GetKtsByAuthor(ctx context.Context, id uuid.UUID) ([]domain.KT, error) {
	res, err := a.client.GetKtsByAuthor(ctx, &pb.GetKtsByAuthorIn{Author: id.String()})
	if err != nil {
		return nil, err
	}

	return mappers.Kt{}.SliceDomain(res.Kts), nil
}

func (a *adapter) UpdateKt(ctx context.Context, in UpdateKtIn) (domain.KT, error) {
	res, err := a.client.UpdateKt(ctx, &pb.UpdateKtIn{
		Id:                 in.Id.String(),
		Checked:            in.Checked,
		ClassProbabilities: mappers.MapToProbabilities(in.ClassProbabilities),
	})
	if err != nil {
		return domain.KT{}, err
	}

	return mappers.Kt{}.Domain(res.Kt), nil
}

func (a *adapter) DeleteKt(ctx context.Context, id uuid.UUID) error {
	_, err := a.client.DeleteKt(ctx, &pb.DeleteKtIn{Id: id.String()})
	return err
}


