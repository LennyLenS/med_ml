package mri

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "exam/internal/generated/grpc/service"
	"exam/internal/server/mappers"
)

func (h *handler) GetMriById(ctx context.Context, in *pb.GetMriByIdIn) (*pb.GetMriByIdOut, error) {
	if _, err := uuid.Parse(in.Id); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id is not a valid uuid: %s", err.Error())
	}

	mri, err := h.services.Mri.GetMriByID(ctx, uuid.MustParse(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	out := new(pb.GetMriByIdOut)
	out.Mri = mappers.MriFromDomain(mri)

	return out, nil
}

func (h *handler) GetMrisByExternalId(ctx context.Context, in *pb.GetMrisByExternalIdIn) (*pb.GetMrisByExternalIdOut, error) {
	if _, err := uuid.Parse(in.ExternalId); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "external_id is not a valid uuid: %s", err.Error())
	}

	mris, err := h.services.Mri.GetMrisByExternalID(ctx, uuid.MustParse(in.ExternalId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	out := new(pb.GetMrisByExternalIdOut)
	out.Mris = mappers.SliceMriFromDomain(mris)

	return out, nil
}

func (h *handler) GetMrisByAuthor(ctx context.Context, in *pb.GetMrisByAuthorIn) (*pb.GetMrisByAuthorOut, error) {
	if _, err := uuid.Parse(in.Author); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "author is not a valid uuid: %s", err.Error())
	}

	mris, err := h.services.Mri.GetMrisByAuthor(ctx, uuid.MustParse(in.Author))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	out := new(pb.GetMrisByAuthorOut)
	out.Mris = mappers.SliceMriFromDomain(mris)

	return out, nil
}

func (h *handler) GetEchographicByMriId(ctx context.Context, in *pb.GetEchographicByMriIdIn) (*pb.GetEchographicByMriIdOut, error) {
	if _, err := uuid.Parse(in.MriId); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "mri_id is not a valid uuid: %s", err.Error())
	}

	echographic, err := h.services.Mri.GetMriEchographicsByID(ctx, uuid.MustParse(in.MriId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	out := new(pb.GetEchographicByMriIdOut)
	out.Echographic = mappers.EchographicFromDomain(echographic)

	return out, nil
}
