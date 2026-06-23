package cytology_image_test

import (
	"context"
	"errors"
	"testing"

	"cytology/internal/domain"
	pb "cytology/internal/generated/grpc/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateCytologyImage_InvalidExternalID(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	_, err := h.CreateCytologyImage(context.Background(), &pb.CreateCytologyImageIn{
		ExternalId: "invalid",
		DoctorId:   uuid.New().String(),
		PatientId:  uuid.New().String(),
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateCytologyImage_InvalidDoctorID(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	_, err := h.CreateCytologyImage(context.Background(), &pb.CreateCytologyImageIn{
		ExternalId: uuid.New().String(),
		DoctorId:   "invalid",
		PatientId:  uuid.New().String(),
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateCytologyImage_InvalidPrevID(t *testing.T) {
	prevID := "invalid"
	h := newHandler(&mockCytologyImageService{})

	_, err := h.CreateCytologyImage(context.Background(), &pb.CreateCytologyImageIn{
		ExternalId: uuid.New().String(),
		DoctorId:   uuid.New().String(),
		PatientId:  uuid.New().String(),
		PrevId:     &prevID,
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateCytologyImage_Success(t *testing.T) {
	id := uuid.New()
	h := newHandler(&mockCytologyImageService{createID: id})

	resp, err := h.CreateCytologyImage(context.Background(), &pb.CreateCytologyImageIn{
		ExternalId:       uuid.New().String(),
		DoctorId:         uuid.New().String(),
		PatientId:        uuid.New().String(),
		DiagnosticNumber: 1,
	})

	require.NoError(t, err)
	require.Equal(t, id.String(), resp.Id)
}

func TestCreateCytologyImage_UnprocessableEntity(t *testing.T) {
	h := newHandler(&mockCytologyImageService{createErr: domain.ErrUnprocessableEntity})

	_, err := h.CreateCytologyImage(context.Background(), &pb.CreateCytologyImageIn{
		ExternalId: uuid.New().String(),
		DoctorId:   uuid.New().String(),
		PatientId:  uuid.New().String(),
	})

	require.Error(t, err)
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestCreateCytologyImage_InternalError(t *testing.T) {
	h := newHandler(&mockCytologyImageService{createErr: errors.New("db error")})

	_, err := h.CreateCytologyImage(context.Background(), &pb.CreateCytologyImageIn{
		ExternalId: uuid.New().String(),
		DoctorId:   uuid.New().String(),
		PatientId:  uuid.New().String(),
	})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
