package mriprocessed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/WantBeASleep/med_ml_lib/dbus"
	"github.com/google/uuid"

	pb "exam/internal/generated/dbus/consume/mriprocessed"
	"exam/internal/services"
	"exam/internal/services/node_segment"
)

var ErrInvalidContor = errors.New("invalid contor")

type subscriber struct {
	services *services.Services
}

func New(
	services *services.Services,
) dbus.Consumer[*pb.MriProcessed] {
	return &subscriber{
		services: services,
	}
}

func (h *subscriber) Consume(ctx context.Context, message *pb.MriProcessed) error {
	if _, err := uuid.Parse(message.MriId); err != nil {
		return fmt.Errorf("mri id is not uuid: %s", message.MriId)
	}

	for _, v := range message.NodesWithSegments {
		for _, segment := range v.Segments {
			if _, err := uuid.Parse(segment.ImageId); err != nil {
				return fmt.Errorf("image id is not uuid: %s", segment.ImageId)
			}
		}
	}

	arg := make([]node_segment.CreateNodesWithSegmentsArg, 0, len(message.NodesWithSegments))
	for _, v := range message.NodesWithSegments {
		node := node_segment.CreateNodesWithSegmentsArgNode{
			Knosp012: v.Node.Knosp_012,
			Knosp3:   v.Node.Knosp_3,
			Knosp4:   v.Node.Knosp_4,
		}

		segments := make([]node_segment.CreateNodesWithSegmentsArgSegment, 0, len(v.Segments))
		for _, segment := range v.Segments {
			// contor json parse
			if !json.Valid(segment.Contor) {
				return ErrInvalidContor
			}

			segments = append(segments, node_segment.CreateNodesWithSegmentsArgSegment{
				ImageID:  uuid.MustParse(segment.ImageId),
				Contor:   json.RawMessage(segment.Contor),
				Knosp012: segment.Knosp_012,
				Knosp3:   segment.Knosp_3,
				Knosp4:   segment.Knosp_4,
			})
		}

		arg = append(arg, node_segment.CreateNodesWithSegmentsArg{
			Node:     node,
			Segments: segments,
		})
	}

	return h.services.NodeSegment.SaveProcessedNodesWithSegments(ctx, uuid.MustParse(message.MriId), arg)
}
