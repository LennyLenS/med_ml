//go:build e2e

package cytology_image_test

import (
	"github.com/stretchr/testify/require"

	pb "cytology/internal/generated/grpc/service"
	"cytology/tests/e2e/flow"
)

func (suite *TestSuite) TestUpdateCytologyImage_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	marking := pb.DiagnosticMarking_DIAGNOSTIC_MARKING_P11
	calcitonin := int32(100)

	resp, err := suite.deps.Adapter.UpdateCytologyImage(
		suite.T().Context(),
		&pb.UpdateCytologyImageIn{
			Id:                data.CytologyImageID.String(),
			DiagnosticMarking: &marking,
			Calcitonin:        &calcitonin,
		},
	)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), marking, *resp.CytologyImage.DiagnosticMarking)
	require.Equal(suite.T(), calcitonin, *resp.CytologyImage.Calcitonin)
}
