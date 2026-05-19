//go:build e2e

package cytology_image_test

import (
	"github.com/stretchr/testify/require"

	pb "cytology/internal/generated/grpc/service"
	"cytology/tests/e2e/flow"
)

func (suite *TestSuite) TestDeleteCytologyImage_Success() {
	data, err := flow.New(suite.deps, flow.CytologyImageInit).Do(suite.T().Context())
	require.NoError(suite.T(), err)

	_, err = suite.deps.Adapter.DeleteCytologyImage(
		suite.T().Context(),
		&pb.DeleteCytologyImageIn{Id: data.CytologyImageID.String()},
	)
	require.NoError(suite.T(), err)

	_, err = suite.deps.Adapter.GetCytologyImageById(
		suite.T().Context(),
		&pb.GetCytologyImageByIdIn{Id: data.CytologyImageID.String()},
	)
	require.Error(suite.T(), err)
}
