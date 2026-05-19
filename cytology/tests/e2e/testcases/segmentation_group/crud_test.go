//go:build e2e

package segmentation_group_test

import (
	"github.com/stretchr/testify/require"

	pb "cytology/internal/generated/grpc/service"
	"cytology/tests/e2e/flow"
)

func (suite *TestSuite) TestCreateSegmentationGroup_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit, flow.SegmentationGroupInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)
	require.NotZero(suite.T(), data.SegmentationGroupID)
}

func (suite *TestSuite) TestGetSegmentationGroupsByCytologyId_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit, flow.SegmentationGroupInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	resp, err := suite.deps.Adapter.GetSegmentationGroupsByCytologyId(
		suite.T().Context(),
		&pb.GetSegmentationGroupsByCytologyIdIn{CytologyId: data.CytologyImageID.String()},
	)
	require.NoError(suite.T(), err)
	require.Len(suite.T(), resp.SegmentationGroups, 1)
	require.Equal(suite.T(), data.SegmentationGroupID, resp.SegmentationGroups[0].Id)
	require.Equal(suite.T(), pb.SegType_SEG_TYPE_NIL, resp.SegmentationGroups[0].SegType)
}

func (suite *TestSuite) TestUpdateSegmentationGroup_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit, flow.SegmentationGroupInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	details := `{"updated":true}`

	resp, err := suite.deps.Adapter.UpdateSegmentationGroup(
		suite.T().Context(),
		&pb.UpdateSegmentationGroupIn{
			Id:      data.SegmentationGroupID,
			Details: &details,
		},
	)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), details, *resp.SegmentationGroup.Details)
	require.Equal(suite.T(), pb.SegType_SEG_TYPE_NIL, resp.SegmentationGroup.SegType)
}

func (suite *TestSuite) TestDeleteSegmentationGroup_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit, flow.SegmentationGroupInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	_, err = suite.deps.Adapter.DeleteSegmentationGroup(
		suite.T().Context(),
		&pb.DeleteSegmentationGroupIn{Id: data.SegmentationGroupID},
	)
	require.NoError(suite.T(), err)

	resp, err := suite.deps.Adapter.GetSegmentationGroupsByCytologyId(
		suite.T().Context(),
		&pb.GetSegmentationGroupsByCytologyIdIn{CytologyId: data.CytologyImageID.String()},
	)
	require.NoError(suite.T(), err)
	require.Empty(suite.T(), resp.SegmentationGroups)
}
