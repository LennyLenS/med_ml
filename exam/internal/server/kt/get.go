package kt

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "exam/internal/generated/grpc/service"
	"exam/internal/server/mappers"
)

func (h *handler) GetKtById(ctx context.Context, in *pb.GetKtByIdIn) (*pb.GetKtByIdOut, error) {
	if _, err := uuid.Parse(in.Id); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id is not a valid uuid: %s", err.Error())
	}

	kt, err := h.services.Kt.GetKtByID(ctx, uuid.MustParse(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	out := new(pb.GetKtByIdOut)
	out.Kt = mappers.KtFromDomain(kt)

	return out, nil
}

func (h *handler) GetKtsByAuthor(ctx context.Context, in *pb.GetKtsByAuthorIn) (*pb.GetKtsByAuthorOut, error) {
	if _, err := uuid.Parse(in.Author); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "author is not a valid uuid: %s", err.Error())
	}

	kts, err := h.services.Kt.GetKtsByAuthor(ctx, uuid.MustParse(in.Author))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	out := new(pb.GetKtsByAuthorOut)
	out.Kts = mappers.SliceKtFromDomain(kts)

	return out, nil
}
