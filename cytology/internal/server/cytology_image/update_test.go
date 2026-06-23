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

func TestUpdateCytologyImage_InvalidID(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	_, err := h.UpdateCytologyImage(context.Background(), &pb.UpdateCytologyImageIn{Id: "invalid"})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestUpdateCytologyImage_NotFound(t *testing.T) {
	h := newHandler(&mockCytologyImageService{updateErr: domain.ErrNotFound})

	_, err := h.UpdateCytologyImage(context.Background(), &pb.UpdateCytologyImageIn{Id: uuid.New().String()})

	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestUpdateCytologyImage_UnprocessableEntity(t *testing.T) {
	h := newHandler(&mockCytologyImageService{updateErr: domain.ErrUnprocessableEntity})

	_, err := h.UpdateCytologyImage(context.Background(), &pb.UpdateCytologyImageIn{Id: uuid.New().String()})

	require.Error(t, err)
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestUpdateCytologyImage_Success(t *testing.T) {
	imageID := uuid.New()
	image := domain.CytologyImage{
		Id:               imageID,
		ExternalID:       uuid.New(),
		DoctorID:         uuid.New(),
		PatientID:        uuid.New(),
		DiagnosticNumber: 3,
		DiagnosDate:      time.Now().UTC(),
		IsLast:           true,
		CreateAt:         time.Now().UTC(),
	}
	h := newHandler(&mockCytologyImageService{updateImage: image})

	resp, err := h.UpdateCytologyImage(context.Background(), &pb.UpdateCytologyImageIn{Id: imageID.String()})

	require.NoError(t, err)
	require.Equal(t, imageID.String(), resp.CytologyImage.Id)
}

func TestUpdateCytologyImage_InternalError(t *testing.T) {
	h := newHandler(&mockCytologyImageService{updateErr: errors.New("db error")})

	_, err := h.UpdateCytologyImage(context.Background(), &pb.UpdateCytologyImageIn{Id: uuid.New().String()})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
