package exam

import (
	"context"

	"github.com/google/uuid"

	"composition-api/internal/adapters/exam/mappers"
	domain "composition-api/internal/domain/exam"
	pb "composition-api/internal/generated/grpc/clients/exam"
)

var mriProjectionMap = map[domain.MriProjection]pb.MriProjection{
	domain.MriProjectionCross: pb.MriProjection_MRI_PROJECTION_CROSS,
	domain.MriProjectionLong:  pb.MriProjection_MRI_PROJECTION_LONG,
}

func (a *adapter) CreateMri(ctx context.Context, in CreateMriIn) (uuid.UUID, error) {
	res, err := a.client.CreateMri(ctx, &pb.CreateMriIn{
		Projection:  mriProjectionMap[in.Projection],
		ExternalId:  in.ExternalID.String(),
		Author:      in.Author.String(),
		DeviceId:    int64(in.DeviceID),
		Description: in.Description,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.MustParse(res.Id), nil
}

func (a *adapter) GetMriById(ctx context.Context, id uuid.UUID) (domain.Mri, error) {
	res, err := a.client.GetMriById(ctx, &pb.GetMriByIdIn{Id: id.String()})
	if err != nil {
		return domain.Mri{}, err
	}

	return mappers.Mri{}.Domain(res.Mri), nil
}

func (a *adapter) GetMrisByExternalId(ctx context.Context, id uuid.UUID) ([]domain.Mri, error) {
	res, err := a.client.GetMrisByExternalId(ctx, &pb.GetMrisByExternalIdIn{ExternalId: id.String()})
	if err != nil {
		return nil, err
	}

	return mappers.Mri{}.SliceDomain(res.Mris), nil
}

func (a *adapter) GetMrisByAuthor(ctx context.Context, id uuid.UUID) ([]domain.Mri, error) {
	res, err := a.client.GetMrisByAuthor(ctx, &pb.GetMrisByAuthorIn{Author: id.String()})
	if err != nil {
		return nil, err
	}

	return mappers.Mri{}.SliceDomain(res.Mris), nil
}

func (a *adapter) GetEchographicByMriId(ctx context.Context, id uuid.UUID) (domain.Echographic, error) {
	res, err := a.client.GetEchographicByMriId(ctx, &pb.GetEchographicByMriIdIn{MriId: id.String()})
	if err != nil {
		return domain.Echographic{}, err
	}

	return mappers.Echographic{}.Domain(res.Echographic), nil
}

func (a *adapter) UpdateMri(ctx context.Context, in UpdateMriIn) (domain.Mri, error) {
	res, err := a.client.UpdateMri(ctx, &pb.UpdateMriIn{
		Id:         in.Id.String(),
		Projection: mappers.PointerFromMap(mriProjectionMap, in.Projection),
		Checked:    in.Checked,
	})
	if err != nil {
		return domain.Mri{}, err
	}

	return mappers.Mri{}.Domain(res.Mri), nil
}

func (a *adapter) UpdateEchographic(ctx context.Context, in domain.Echographic) (domain.Echographic, error) {
	res, err := a.client.UpdateEchographic(ctx, &pb.UpdateEchographicIn{
		Echographic: &pb.Echographic{
			Id:              in.Id.String(),
			Contors:         in.Contors,
			LeftLobeLength:  in.LeftLobeLength,
			LeftLobeWidth:   in.LeftLobeWidth,
			LeftLobeThick:   in.LeftLobeThick,
			LeftLobeVolum:   in.LeftLobeVolum,
			RightLobeLength: in.RightLobeLength,
			RightLobeWidth:  in.RightLobeWidth,
			RightLobeThick:  in.RightLobeThick,
			RightLobeVolum:  in.RightLobeVolum,
			GlandVolum:      in.GlandVolum,
			Isthmus:         in.Isthmus,
			Struct:          in.Struct,
			Echogenicity:    in.Echogenicity,
			RegionalLymph:   in.RegionalLymph,
			Vascularization: in.Vascularization,
			Location:        in.Location,
			Additional:      in.Additional,
			Conclusion:      in.Conclusion,
		},
	})
	if err != nil {
		return domain.Echographic{}, err
	}

	return mappers.Echographic{}.Domain(res.Echographic), nil
}

func (a *adapter) DeleteMri(ctx context.Context, id uuid.UUID) error {
	_, err := a.client.DeleteMri(ctx, &pb.DeleteMriIn{Id: id.String()})
	return err
}


