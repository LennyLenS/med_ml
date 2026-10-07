package kt

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "exam/internal/generated/grpc/service"
	"exam/internal/server/mappers"
	"exam/internal/services/kt"
)

func (h *handler) UpdateKt(ctx context.Context, in *pb.UpdateKtIn) (*pb.UpdateKtOut, error) {
	if _, err := uuid.Parse(in.Id); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id is not a valid uuid: %s", err.Error())
	}

	var classProbabilities *map[string]float64
	if in.ClassProbabilities != nil {
		classProbabilities = &in.ClassProbabilities.ClassProbabilities
	}

	kt, err := h.services.Kt.UpdateKt(ctx, kt.UpdateKTArg{
		Id:               uuid.MustParse(in.Id),
		Checked:          in.Checked,
		PredictedClasses: classProbabilities,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Что то пошло не так: %s", err.Error())
	}

	out := new(pb.UpdateKtOut)
	out.Kt = mappers.KtFromDomain(kt)

	return out, nil
}
