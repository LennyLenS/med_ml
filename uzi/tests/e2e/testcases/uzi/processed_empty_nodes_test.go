//go:build e2e

package uzi_test

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	pbDbus "uzi/internal/generated/dbus/consume/uziprocessed"
	pb "uzi/internal/generated/grpc/service"
	"uzi/tests/e2e/flow"
)

func (suite *TestSuite) TestUziProcessed_EmptyNodes() {
	data, err := flow.New(suite.deps, flow.DeviceInit, flow.UziInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	request := &pbDbus.UziProcessed{
		UziId: data.Uzi.Id.String(),
		//специально делаем nodes_with_segments пустыми
	}

	message, err := proto.Marshal(request)
	require.NoError(suite.T(), err)

	_, _, err = suite.deps.Dbus.SendMessage(
		&sarama.ProducerMessage{
			Topic: "uziprocessed",
			Value: sarama.ByteEncoder(message),
		},
	)
	require.NoError(suite.T(), err)

	ctx, cancel := context.WithTimeout(suite.T().Context(), 20*time.Second)
	defer cancel()

	backoff := time.Second
	for {
		select {
		case <-ctx.Done():
			require.FailNow(suite.T(), "context done. uzi not completed")
		case <-time.After(backoff):
			resp, err := suite.deps.Adapter.GetUziById(ctx, &pb.GetUziByIdIn{Id: data.Uzi.Id.String()})
			require.NoError(suite.T(), err)
			if resp.Uzi.Status != pb.UziStatus_UZI_STATUS_COMPLETED {
				backoff *= 2
				continue
			}

			nodesResp, err := suite.deps.Adapter.GetNodesByUziId(ctx, &pb.GetNodesByUziIdIn{UziId: data.Uzi.Id.String()})
			if err != nil {
				require.ErrorContains(suite.T(), err, "not found")
				return
			}
			require.Len(suite.T(), nodesResp.Nodes, 0)
			return
		}
	}
}
