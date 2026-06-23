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

func TestUpdateOriginalImage_InvalidID(t *testing.T) {
	h := newHandler(&mockOriginalImageService{})

	_, err := h.UpdateOriginalImage(context.Background(), &pb.UpdateOriginalImageIn{Id: "invalid"})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestUpdateOriginalImage_Success(t *testing.T) {
	imageID := uuid.New()
	delay := 1.5
	viewed := true
	h := newHandler(&mockOriginalImageService{
		updateImage: domain.OriginalImage{
			Id:         imageID,
			CytologyID: uuid.New(),
			ImagePath:  "path",
			CreateDate: time.Now().UTC(),
			DelayTime:  &delay,
			ViewedFlag: viewed,
		},
	})

	resp, err := h.UpdateOriginalImage(context.Background(), &pb.UpdateOriginalImageIn{
		Id:         imageID.String(),
		DelayTime:  &delay,
		ViewedFlag: &viewed,
	})

	require.NoError(t, err)
	require.Equal(t, imageID.String(), resp.OriginalImage.Id)
	require.True(t, resp.OriginalImage.ViewedFlag)
}

func TestUpdateOriginalImage_InternalError(t *testing.T) {
	h := newHandler(&mockOriginalImageService{updateErr: errors.New("db error")})

	_, err := h.UpdateOriginalImage(context.Background(), &pb.UpdateOriginalImageIn{Id: uuid.New().String()})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
