//go:build e2e

package segmentation_test

import (
	"github.com/stretchr/testify/require"

	pb "cytology/internal/generated/grpc/service"
	"cytology/tests/e2e/flow"
)

func (suite *TestSuite) TestCreateSegmentation_Success() {
	data, err := flow.New(
		suite.deps,
		flow.CytologyImageInit,
		flow.SegmentationGroupInit,
		flow.SegmentationInit,
	).Do(suite.T().Context())
	require.NoError(suite.T(), err)
	require.NotZero(suite.T(), data.SegmentationID)
}

func (suite *TestSuite) TestGetSegmentationById_Success() {
	data, err := flow.New(
		suite.deps,
		flow.CytologyImageInit,
		flow.SegmentationGroupInit,
		flow.SegmentationInit,
	).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	resp, err := suite.deps.Adapter.GetSegmentationById(
		suite.T().Context(),
		&pb.GetSegmentationByIdIn{Id: data.SegmentationID},
	)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), data.SegmentationID, resp.Segmentation.Id)
	require.Len(suite.T(), resp.Segmentation.Points, 2)
}

func (suite *TestSuite) TestGetSegmentsByGroupId_Success() {
	data, err := flow.New(
		suite.deps,
		flow.CytologyImageInit,
		flow.SegmentationGroupInit,
		flow.SegmentationInit,
	).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	resp, err := suite.deps.Adapter.GetSegmentsByGroupId(
		suite.T().Context(),
		&pb.GetSegmentsByGroupIdIn{SegmentationGroupId: data.SegmentationGroupID},
	)
	require.NoError(suite.T(), err)
	require.Len(suite.T(), resp.Segmentations, 1)
	require.Equal(suite.T(), data.SegmentationID, resp.Segmentations[0].Id)
}

func (suite *TestSuite) TestUpdateSegmentation_Success() {
	data, err := flow.New(
		suite.deps,
		flow.CytologyImageInit,
		flow.SegmentationGroupInit,
		flow.SegmentationInit,
	).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	resp, err := suite.deps.Adapter.UpdateSegmentation(
		suite.T().Context(),
		&pb.UpdateSegmentationIn{
			Id: data.SegmentationID,
			Points: []*pb.SegmentationPointCreate{
				{X: 1, Y: 2},
				{X: 3, Y: 4},
				{X: 5, Y: 6},
			},
		},
	)
	require.NoError(suite.T(), err)
	require.Len(suite.T(), resp.Segmentation.Points, 3)
}

func (suite *TestSuite) TestDeleteSegmentation_Success() {
	data, err := flow.New(
		suite.deps,
		flow.CytologyImageInit,
		flow.SegmentationGroupInit,
		flow.SegmentationInit,
	).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	_, err = suite.deps.Adapter.DeleteSegmentation(
		suite.T().Context(),
		&pb.DeleteSegmentationIn{Id: data.SegmentationID},
	)
	require.NoError(suite.T(), err)

	resp, err := suite.deps.Adapter.GetSegmentsByGroupId(
		suite.T().Context(),
		&pb.GetSegmentsByGroupIdIn{SegmentationGroupId: data.SegmentationGroupID},
	)
	require.NoError(suite.T(), err)
	require.Empty(suite.T(), resp.Segmentations)
}
