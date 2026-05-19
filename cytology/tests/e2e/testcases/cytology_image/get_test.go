//go:build e2e

package cytology_image_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pb "cytology/internal/generated/grpc/service"
	"cytology/tests/e2e/flow"
)

func (suite *TestSuite) TestGetCytologyImageById_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	resp, err := suite.deps.Adapter.GetCytologyImageById(
		suite.T().Context(),
		&pb.GetCytologyImageByIdIn{Id: data.CytologyImageID.String()},
	)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), data.CytologyImageID.String(), resp.CytologyImage.Id)
}

func (suite *TestSuite) TestGetCytologyImagesByExternalId_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	resp, err := suite.deps.Adapter.GetCytologyImagesByExternalId(
		suite.T().Context(),
		&pb.GetCytologyImagesByExternalIdIn{ExternalId: data.ExternalID.String()},
	)
	require.NoError(suite.T(), err)
	require.Len(suite.T(), resp.CytologyImages, 1)
	require.Equal(suite.T(), data.CytologyImageID.String(), resp.CytologyImages[0].Id)
}

func (suite *TestSuite) TestGetCytologyImagesByDoctorIdAndPatientId_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	resp, err := suite.deps.Adapter.GetCytologyImagesByDoctorIdAndPatientId(
		suite.T().Context(),
		&pb.GetCytologyImagesByDoctorIdAndPatientIdIn{
			DoctorId:  data.DoctorID.String(),
			PatientId: data.PatientID.String(),
		},
	)
	require.NoError(suite.T(), err)
	require.Len(suite.T(), resp.CytologyImages, 1)
	require.Equal(suite.T(), data.CytologyImageID.String(), resp.CytologyImages[0].Id)
}

func (suite *TestSuite) TestGetCytologyImagesByExternalId_Empty() {
	resp, err := suite.deps.Adapter.GetCytologyImagesByExternalId(
		suite.T().Context(),
		&pb.GetCytologyImagesByExternalIdIn{ExternalId: uuid.New().String()},
	)
	require.NoError(suite.T(), err)
	require.Empty(suite.T(), resp.CytologyImages)
}
