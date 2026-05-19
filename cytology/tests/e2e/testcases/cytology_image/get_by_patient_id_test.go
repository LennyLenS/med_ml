//go:build e2e

package cytology_image_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pb "cytology/internal/generated/grpc/service"
	"cytology/tests/e2e/flow"
)

func (suite *TestSuite) TestGetCytologyImagesByPatientId_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	resp, err := suite.deps.Adapter.GetCytologyImagesByPatientId(
		suite.T().Context(),
		&pb.GetCytologyImagesByPatientIdIn{PatientId: data.PatientID.String()},
	)
	require.NoError(suite.T(), err)
	require.Len(suite.T(), resp.CytologyImages, 1)
	require.Equal(suite.T(), data.CytologyImageID.String(), resp.CytologyImages[0].Id)
}

func (suite *TestSuite) TestGetCytologyImagesByPatientId_Empty() {
	resp, err := suite.deps.Adapter.GetCytologyImagesByPatientId(
		suite.T().Context(),
		&pb.GetCytologyImagesByPatientIdIn{PatientId: uuid.New().String()},
	)
	require.NoError(suite.T(), err)
	require.Empty(suite.T(), resp.CytologyImages)
}
