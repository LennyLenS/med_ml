package original_image_test

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

func TestGetOriginalImageById_InvalidID(t *testing.T) {
	h := newHandler(&mockOriginalImageService{})

	_, err := h.GetOriginalImageById(context.Background(), &pb.GetOriginalImageByIdIn{Id: "invalid"})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetOriginalImageById_Success(t *testing.T) {
	imageID := uuid.New()
	cytologyID := uuid.New()
	h := newHandler(&mockOriginalImageService{
		image: domain.OriginalImage{
			Id:         imageID,
			CytologyID: cytologyID,
			ImagePath:  "path/to/image",
			CreateDate: time.Now().UTC(),
			ViewedFlag: true,
		},
	})

	resp, err := h.GetOriginalImageById(context.Background(), &pb.GetOriginalImageByIdIn{Id: imageID.String()})

	require.NoError(t, err)
	require.Equal(t, imageID.String(), resp.OriginalImage.Id)
	require.Equal(t, cytologyID.String(), resp.OriginalImage.CytologyId)
}

func TestGetOriginalImageById_InternalError(t *testing.T) {
	h := newHandler(&mockOriginalImageService{imageErr: errors.New("db error")})

	_, err := h.GetOriginalImageById(context.Background(), &pb.GetOriginalImageByIdIn{Id: uuid.New().String()})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}

func TestGetOriginalImagesByCytologyId_InvalidCytologyID(t *testing.T) {
	h := newHandler(&mockOriginalImageService{})

	_, err := h.GetOriginalImagesByCytologyId(context.Background(), &pb.GetOriginalImagesByCytologyIdIn{
		CytologyId: "invalid",
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetOriginalImagesByCytologyId_EmptyOnError(t *testing.T) {
	h := newHandler(&mockOriginalImageService{imagesErr: errors.New("db error")})

	resp, err := h.GetOriginalImagesByCytologyId(context.Background(), &pb.GetOriginalImagesByCytologyIdIn{
		CytologyId: uuid.New().String(),
	})

	require.NoError(t, err)
	require.Empty(t, resp.OriginalImages)
}

func TestGetOriginalImagesByCytologyId_Success(t *testing.T) {
	imageID := uuid.New()
	h := newHandler(&mockOriginalImageService{
		images: []domain.OriginalImage{
			{
				Id:         imageID,
				CytologyID: uuid.New(),
				ImagePath:  "path",
				CreateDate: time.Now().UTC(),
			},
		},
	})

	resp, err := h.GetOriginalImagesByCytologyId(context.Background(), &pb.GetOriginalImagesByCytologyIdIn{
		CytologyId: uuid.New().String(),
	})

	require.NoError(t, err)
	require.Len(t, resp.OriginalImages, 1)
	require.Equal(t, imageID.String(), resp.OriginalImages[0].Id)
}
