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

func TestDeleteCytologyImage_InvalidID(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	_, err := h.DeleteCytologyImage(context.Background(), &pb.DeleteCytologyImageIn{Id: "invalid"})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestDeleteCytologyImage_NotFound(t *testing.T) {
	h := newHandler(&mockCytologyImageService{deleteErr: domain.ErrNotFound})

	_, err := h.DeleteCytologyImage(context.Background(), &pb.DeleteCytologyImageIn{Id: uuid.New().String()})

	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestDeleteCytologyImage_Success(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	resp, err := h.DeleteCytologyImage(context.Background(), &pb.DeleteCytologyImageIn{Id: uuid.New().String()})

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestDeleteCytologyImage_InternalError(t *testing.T) {
	h := newHandler(&mockCytologyImageService{deleteErr: errors.New("db error")})

	_, err := h.DeleteCytologyImage(context.Background(), &pb.DeleteCytologyImageIn{Id: uuid.New().String()})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
