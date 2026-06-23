package original_image_test

import (
	"context"
	"errors"
	"testing"

	pb "cytology/internal/generated/grpc/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateOriginalImage_InvalidCytologyID(t *testing.T) {
	h := newHandler(&mockOriginalImageService{})

	_, err := h.CreateOriginalImage(context.Background(), &pb.CreateOriginalImageIn{
		CytologyId: "invalid",
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateOriginalImage_MissingFileAndPath(t *testing.T) {
	h := newHandler(&mockOriginalImageService{})

	_, err := h.CreateOriginalImage(context.Background(), &pb.CreateOriginalImageIn{
		CytologyId: uuid.New().String(),
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateOriginalImage_MissingContentType(t *testing.T) {
	h := newHandler(&mockOriginalImageService{})

	_, err := h.CreateOriginalImage(context.Background(), &pb.CreateOriginalImageIn{
		CytologyId: uuid.New().String(),
		File:       []byte("data"),
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateOriginalImage_WithImagePath_Success(t *testing.T) {
	id := uuid.New()
	h := newHandler(&mockOriginalImageService{createID: id})

	imagePath := "cytology/image/path"
	resp, err := h.CreateOriginalImage(context.Background(), &pb.CreateOriginalImageIn{
		CytologyId: uuid.New().String(),
		ImagePath:  &imagePath,
	})

	require.NoError(t, err)
	require.Equal(t, id.String(), resp.Id)
}

func TestCreateOriginalImage_WithFile_Success(t *testing.T) {
	id := uuid.New()
	contentType := "image/png"
	h := newHandler(&mockOriginalImageService{createID: id})

	resp, err := h.CreateOriginalImage(context.Background(), &pb.CreateOriginalImageIn{
		CytologyId:  uuid.New().String(),
		File:        []byte("data"),
		ContentType: contentType,
	})

	require.NoError(t, err)
	require.Equal(t, id.String(), resp.Id)
}

func TestCreateOriginalImage_InternalError(t *testing.T) {
	h := newHandler(&mockOriginalImageService{createErr: errors.New("db error")})

	imagePath := "path"
	_, err := h.CreateOriginalImage(context.Background(), &pb.CreateOriginalImageIn{
		CytologyId: uuid.New().String(),
		ImagePath:  &imagePath,
	})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
