//go:build e2e

package cytology_image_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pb "cytology/internal/generated/grpc/service"
	"cytology/tests/e2e/flow"
)

func (suite *TestSuite) TestGetCytologyImageIdsByDoctorIdAndPatientId_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	resp, err := suite.deps.Adapter.GetCytologyImageIdsByDoctorIdAndPatientId(
		suite.T().Context(),
		&pb.GetCytologyImageIdsByDoctorIdAndPatientIdIn{
			DoctorId:  data.DoctorID.String(),
			PatientId: data.PatientID.String(),
		},
	)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), []string{data.CytologyImageID.String()}, resp.Ids)
}

func (suite *TestSuite) TestGetCytologyImageIdsByDoctorIdAndPatientId_Empty() {
	resp, err := suite.deps.Adapter.GetCytologyImageIdsByDoctorIdAndPatientId(
		suite.T().Context(),
		&pb.GetCytologyImageIdsByDoctorIdAndPatientIdIn{
			DoctorId:  uuid.New().String(),
			PatientId: uuid.New().String(),
		},
	)
	require.NoError(suite.T(), err)
	require.Empty(suite.T(), resp.Ids)
}

func (suite *TestSuite) TestGetCytologyImageIdsByDoctorIdAndPatientId_InvalidDoctorID() {
	_, err := suite.deps.Adapter.GetCytologyImageIdsByDoctorIdAndPatientId(
		suite.T().Context(),
		&pb.GetCytologyImageIdsByDoctorIdAndPatientIdIn{
			DoctorId:  "invalid",
			PatientId: uuid.New().String(),
		},
	)
	require.Error(suite.T(), err)
}
