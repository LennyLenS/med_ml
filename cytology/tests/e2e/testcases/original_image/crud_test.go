//go:build e2e

package original_image_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pb "cytology/internal/generated/grpc/service"
	"cytology/tests/e2e/flow"
)

func (suite *TestSuite) TestCreateOriginalImage_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit, flow.OriginalImageInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)
	require.NotEqual(suite.T(), uuid.Nil, data.OriginalImageID)
}

func (suite *TestSuite) TestGetOriginalImageById_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit, flow.OriginalImageInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	resp, err := suite.deps.Adapter.GetOriginalImageById(
		suite.T().Context(),
		&pb.GetOriginalImageByIdIn{Id: data.OriginalImageID.String()},
	)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), data.OriginalImageID.String(), resp.OriginalImage.Id)
	require.Equal(suite.T(), data.CytologyImageID.String(), resp.OriginalImage.CytologyId)
}

func (suite *TestSuite) TestGetOriginalImagesByCytologyId_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit, flow.OriginalImageInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	resp, err := suite.deps.Adapter.GetOriginalImagesByCytologyId(
		suite.T().Context(),
		&pb.GetOriginalImagesByCytologyIdIn{CytologyId: data.CytologyImageID.String()},
	)
	require.NoError(suite.T(), err)
	require.Len(suite.T(), resp.OriginalImages, 1)
	require.Equal(suite.T(), data.OriginalImageID.String(), resp.OriginalImages[0].Id)
}

func (suite *TestSuite) TestUpdateOriginalImage_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit, flow.OriginalImageInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	delay := 1.5
	viewed := true

	resp, err := suite.deps.Adapter.UpdateOriginalImage(
		suite.T().Context(),
		&pb.UpdateOriginalImageIn{
			Id:         data.OriginalImageID.String(),
			DelayTime:  &delay,
			ViewedFlag: &viewed,
		},
	)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), delay, *resp.OriginalImage.DelayTime)
	require.True(suite.T(), resp.OriginalImage.ViewedFlag)
}

func (suite *TestSuite) TestCreateOriginalImage_InvalidCytologyId() {
	_, err := suite.deps.Adapter.CreateOriginalImage(
		suite.T().Context(),
		&pb.CreateOriginalImageIn{
			CytologyId: "invalid",
			ImagePath:  ptr("path"),
		},
	)
	require.Error(suite.T(), err)
}

func ptr(s string) *string {
	return &s
}
