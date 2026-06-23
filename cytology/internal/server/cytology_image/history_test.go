package cytology_image_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"cytology/internal/domain"
	pb "cytology/internal/generated/grpc/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetCytologyImageHistory_InvalidID(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	_, err := h.GetCytologyImageHistory(context.Background(), &pb.GetCytologyImageHistoryIn{Id: "invalid"})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetCytologyImageHistory_Success(t *testing.T) {
	imageID := uuid.New()
	h := newHandler(&mockCytologyImageService{
		history: []domain.CytologyImage{
			{
				Id:               imageID,
				ExternalID:       uuid.New(),
				DoctorID:         uuid.New(),
				PatientID:        uuid.New(),
				DiagnosticNumber: 1,
				DiagnosDate:      time.Now().UTC(),
				IsLast:           false,
				CreateAt:         time.Now().UTC(),
			},
			{
				Id:               uuid.New(),
				ExternalID:       uuid.New(),
				DoctorID:         uuid.New(),
				PatientID:        uuid.New(),
				DiagnosticNumber: 2,
				DiagnosDate:      time.Now().UTC(),
				IsLast:           true,
				CreateAt:         time.Now().UTC(),
			},
		},
	})

	resp, err := h.GetCytologyImageHistory(context.Background(), &pb.GetCytologyImageHistoryIn{Id: imageID.String()})

	require.NoError(t, err)
	require.Len(t, resp.CytologyImages, 2)
}

func TestGetCytologyImageHistory_InternalError(t *testing.T) {
	h := newHandler(&mockCytologyImageService{historyErr: errors.New("db error")})

	_, err := h.GetCytologyImageHistory(context.Background(), &pb.GetCytologyImageHistoryIn{Id: uuid.New().String()})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
