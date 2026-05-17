package flow

import (
	"context"
	"fmt"

	pb "cytology/internal/generated/grpc/service"
)

var SegmentationGroupInit flowfuncDepsInjector = func(deps *Deps) flowfunc {
	return func(ctx context.Context, data FlowData) (FlowData, error) {
		resp, err := deps.Adapter.CreateSegmentationGroup(ctx, &pb.CreateSegmentationGroupIn{
			CytologyId: data.CytologyImageID.String(),
			SegType:    pb.SegType_SEG_TYPE_NIL,
			GroupType:  pb.GroupType_GROUP_TYPE_CE,
			IsAi:       false,
		})
		if err != nil {
			return FlowData{}, fmt.Errorf("create segmentation group: %w", err)
		}

		data.SegmentationGroupID = resp.Id
		return data, nil
	}
}
