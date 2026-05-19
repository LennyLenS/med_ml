//go:build e2e

package cytology_image_test

import (
	"github.com/stretchr/testify/require"

	pb "cytology/internal/generated/grpc/service"
	"cytology/tests/e2e/flow"
)

func (suite *TestSuite) TestGetCytologyImageHistory_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	copyResp, err := suite.deps.Adapter.CopyCytologyImage(
		suite.T().Context(),
		&pb.CopyCytologyImageIn{Id: data.CytologyImageID.String()},
	)
	require.NoError(suite.T(), err)

	historyResp, err := suite.deps.Adapter.GetCytologyImageHistory(
		suite.T().Context(),
		&pb.GetCytologyImageHistoryIn{Id: copyResp.CytologyImage.Id},
	)
	require.NoError(suite.T(), err)
	require.Len(suite.T(), historyResp.CytologyImages, 2)

	ids := []string{
		historyResp.CytologyImages[0].Id,
		historyResp.CytologyImages[1].Id,
	}
	require.Contains(suite.T(), ids, data.CytologyImageID.String())
	require.Contains(suite.T(), ids, copyResp.CytologyImage.Id)
}
