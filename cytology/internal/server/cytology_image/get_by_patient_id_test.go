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

func TestGetCytologyImagesByPatientId_InvalidPatientID(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	_, err := h.GetCytologyImagesByPatientId(context.Background(), &pb.GetCytologyImagesByPatientIdIn{
		PatientId: "invalid",
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetCytologyImagesByPatientId_Empty(t *testing.T) {
	h := newHandler(&mockCytologyImageService{images: nil})

	resp, err := h.GetCytologyImagesByPatientId(context.Background(), &pb.GetCytologyImagesByPatientIdIn{
		PatientId: uuid.New().String(),
	})

	require.NoError(t, err)
	require.Empty(t, resp.CytologyImages)
}

func TestGetCytologyImagesByPatientId_Success(t *testing.T) {
	imageID := uuid.New()
	patientID := uuid.New()
	marking := domain.DiagnosticMarkingP11
	images := []domain.CytologyImage{
		{
			Id:                imageID,
			PatientID:         patientID,
			DiagnosticNumber:  1,
			DiagnosticMarking: &marking,
			DiagnosDate:       time.Now().UTC(),
			IsLast:            true,
			CreateAt:          time.Now().UTC(),
		},
	}

	h := newHandler(&mockCytologyImageService{images: images})

	resp, err := h.GetCytologyImagesByPatientId(context.Background(), &pb.GetCytologyImagesByPatientIdIn{
		PatientId: patientID.String(),
	})

	require.NoError(t, err)
	require.Len(t, resp.CytologyImages, 1)
	require.Equal(t, imageID.String(), resp.CytologyImages[0].Id)
}

func TestGetCytologyImagesByPatientId_InternalError(t *testing.T) {
	h := newHandler(&mockCytologyImageService{imagesErr: errors.New("db error")})

	_, err := h.GetCytologyImagesByPatientId(context.Background(), &pb.GetCytologyImagesByPatientIdIn{
		PatientId: uuid.New().String(),
	})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
