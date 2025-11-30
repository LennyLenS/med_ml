package cytology

import (
	"context"
	"io"

	"composition-api/internal/adapters/cytology/mappers"
	adapter_errors "composition-api/internal/adapters/errors"
	domain "composition-api/internal/domain/cytology"
	pb "composition-api/internal/generated/grpc/clients/cytology"

	"github.com/google/uuid"
)

func (a *adapter) CreateOriginalImage(ctx context.Context, in CreateOriginalImageIn) (uuid.UUID, error) {
	fileData, err := io.ReadAll(in.File.File)
	if err != nil {
		return uuid.Nil, adapter_errors.HandleGRPCError(err)
	}

	req := &pb.CreateOriginalImageIn{
		CytologyId:  in.CytologyID.String(),
		File:        fileData,
		ContentType: in.ContentType,
		DelayTime:   in.DelayTime,
	}

	res, err := a.client.CreateOriginalImage(ctx, req)
	if err != nil {
		return uuid.Nil, adapter_errors.HandleGRPCError(err)
	}

	return uuid.MustParse(res.Id), nil
}

func (a *adapter) GetOriginalImageById(ctx context.Context, id uuid.UUID) (domain.OriginalImage, error) {
	res, err := a.client.GetOriginalImageById(ctx, &pb.GetOriginalImageByIdIn{Id: id.String()})
	if err != nil {
		return domain.OriginalImage{}, adapter_errors.HandleGRPCError(err)
	}

	return mappers.OriginalImage{}.Domain(res.OriginalImage), nil
}

func (a *adapter) GetOriginalImagesByCytologyId(ctx context.Context, id uuid.UUID) ([]domain.OriginalImage, error) {
	res, err := a.client.GetOriginalImagesByCytologyId(ctx, &pb.GetOriginalImagesByCytologyIdIn{CytologyId: id.String()})
	if err != nil {
		return nil, adapter_errors.HandleGRPCError(err)
	}

	return mappers.OriginalImage{}.SliceDomain(res.OriginalImages), nil
}

func (a *adapter) UpdateOriginalImage(ctx context.Context, in UpdateOriginalImageIn) (domain.OriginalImage, error) {
	req := &pb.UpdateOriginalImageIn{
		Id:         in.Id.String(),
		DelayTime:  in.DelayTime,
		ViewedFlag: in.ViewedFlag,
	}

	res, err := a.client.UpdateOriginalImage(ctx, req)
	if err != nil {
		return domain.OriginalImage{}, adapter_errors.HandleGRPCError(err)
	}

	return mappers.OriginalImage{}.Domain(res.OriginalImage), nil
}
