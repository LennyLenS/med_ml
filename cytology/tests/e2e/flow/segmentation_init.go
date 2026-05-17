package flow

import (
	"context"
	"fmt"

	pb "cytology/internal/generated/grpc/service"
)

var SegmentationInit flowfuncDepsInjector = func(deps *Deps) flowfunc {
	return func(ctx context.Context, data FlowData) (FlowData, error) {
		resp, err := deps.Adapter.CreateSegmentation(ctx, &pb.CreateSegmentationIn{
			SegmentationGroupId: data.SegmentationGroupID,
			Points: []*pb.SegmentationPointCreate{
				{X: 10, Y: 20},
				{X: 30, Y: 40},
			},
		})
		if err != nil {
			return FlowData{}, fmt.Errorf("create segmentation: %w", err)
		}

		data.SegmentationID = resp.Id
		return data, nil
	}
}
