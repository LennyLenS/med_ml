//go:build e2e

package cytology_image_test

import (
	"github.com/stretchr/testify/require"

	pb "cytology/internal/generated/grpc/service"
	"cytology/tests/e2e/flow"
)

func (suite *TestSuite) TestCopyCytologyImage_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit, flow.SegmentationGroupInit, flow.SegmentationInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	copyResp, err := suite.deps.Adapter.CopyCytologyImage(
		suite.T().Context(),
		&pb.CopyCytologyImageIn{Id: data.CytologyImageID.String()},
	)
	require.NoError(suite.T(), err)
	require.NotEqual(suite.T(), data.CytologyImageID.String(), copyResp.CytologyImage.Id)
	require.True(suite.T(), copyResp.CytologyImage.IsLast)
	require.NotNil(suite.T(), copyResp.CytologyImage.PrevId)
	require.Equal(suite.T(), data.CytologyImageID.String(), *copyResp.CytologyImage.PrevId)

	oldResp, err := suite.deps.Adapter.GetCytologyImageById(
		suite.T().Context(),
		&pb.GetCytologyImageByIdIn{Id: data.CytologyImageID.String()},
	)
	require.NoError(suite.T(), err)
	require.False(suite.T(), oldResp.CytologyImage.IsLast)

	groupsResp, err := suite.deps.Adapter.GetSegmentationGroupsByCytologyId(
		suite.T().Context(),
		&pb.GetSegmentationGroupsByCytologyIdIn{CytologyId: copyResp.CytologyImage.Id},
	)
	require.NoError(suite.T(), err)
	require.Len(suite.T(), groupsResp.SegmentationGroups, 1)
}
