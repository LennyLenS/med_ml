//go:build e2e

package cytology_image_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pb "cytology/internal/generated/grpc/service"
	"cytology/tests/e2e/flow"
)

func (suite *TestSuite) TestCreateCytologyImage_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	getResp, err := suite.deps.Adapter.GetCytologyImageById(
		suite.T().Context(),
		&pb.GetCytologyImageByIdIn{Id: data.CytologyImageID.String()},
	)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), data.CytologyImageID.String(), getResp.CytologyImage.Id)
	require.Equal(suite.T(), data.ExternalID.String(), getResp.CytologyImage.ExternalId)
	require.Equal(suite.T(), data.DoctorID.String(), getResp.CytologyImage.DoctorId)
	require.Equal(suite.T(), data.PatientID.String(), getResp.CytologyImage.PatientId)
	require.Equal(suite.T(), int32(1), getResp.CytologyImage.DiagnosticNumber)
	require.True(suite.T(), getResp.CytologyImage.IsLast)
}

func (suite *TestSuite) TestCreateCytologyImage_InvalidExternalID() {
	_, err := suite.deps.Adapter.CreateCytologyImage(
		suite.T().Context(),
		&pb.CreateCytologyImageIn{
			ExternalId:       "invalid",
			DoctorId:         uuid.New().String(),
			PatientId:        uuid.New().String(),
			DiagnosticNumber: 1,
		},
	)
	require.Error(suite.T(), err)
}
