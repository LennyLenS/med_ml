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

func TestCopyCytologyImage_InvalidID(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	_, err := h.CopyCytologyImage(context.Background(), &pb.CopyCytologyImageIn{Id: "invalid"})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCopyCytologyImage_BadRequest(t *testing.T) {
	h := newHandler(&mockCytologyImageService{copyErr: domain.ErrBadRequest})

	_, err := h.CopyCytologyImage(context.Background(), &pb.CopyCytologyImageIn{Id: uuid.New().String()})

	require.Error(t, err)
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestCopyCytologyImage_Success(t *testing.T) {
	newID := uuid.New()
	image := domain.CytologyImage{
		Id:               newID,
		ExternalID:       uuid.New(),
		DoctorID:         uuid.New(),
		PatientID:        uuid.New(),
		DiagnosticNumber: 1,
		DiagnosDate:      time.Now().UTC(),
		IsLast:           true,
		CreateAt:         time.Now().UTC(),
	}
	h := newHandler(&mockCytologyImageService{copyImage: image})

	resp, err := h.CopyCytologyImage(context.Background(), &pb.CopyCytologyImageIn{Id: uuid.New().String()})

	require.NoError(t, err)
	require.Equal(t, newID.String(), resp.CytologyImage.Id)
}

func TestCopyCytologyImage_InternalError(t *testing.T) {
	h := newHandler(&mockCytologyImageService{copyErr: errors.New("db error")})

	_, err := h.CopyCytologyImage(context.Background(), &pb.CopyCytologyImageIn{Id: uuid.New().String()})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
