package flow

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	pb "cytology/internal/generated/grpc/service"
)

var OriginalImageInit flowfuncDepsInjector = func(deps *Deps) flowfunc {
	return func(ctx context.Context, data FlowData) (FlowData, error) {
		imageID := uuid.New()
		imagePath := fmt.Sprintf("%s/%s/%s", data.CytologyImageID, imageID, imageID)

		resp, err := deps.Adapter.CreateOriginalImage(ctx, &pb.CreateOriginalImageIn{
			CytologyId: data.CytologyImageID.String(),
			ImagePath:  &imagePath,
		})
		if err != nil {
			return FlowData{}, fmt.Errorf("create original image: %w", err)
		}

		data.OriginalImageID = uuid.MustParse(resp.Id)
		return data, nil
	}
}
