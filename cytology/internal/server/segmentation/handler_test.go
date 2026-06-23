package segmentation_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"cytology/internal/domain"
	pb "cytology/internal/generated/grpc/service"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateSegmentation_Success(t *testing.T) {
	h := newHandler(&mockSegmentationService{createID: 42})

	resp, err := h.CreateSegmentation(context.Background(), &pb.CreateSegmentationIn{
		SegmentationGroupId: 1,
		Points: []*pb.SegmentationPointCreate{
			{X: 10, Y: 20},
			{X: 30, Y: 40},
		},
	})

	require.NoError(t, err)
	require.Equal(t, int32(42), resp.Id)
}

func TestCreateSegmentation_InternalError(t *testing.T) {
	h := newHandler(&mockSegmentationService{createErr: errors.New("db error")})

	_, err := h.CreateSegmentation(context.Background(), &pb.CreateSegmentationIn{
		SegmentationGroupId: 1,
	})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}

func TestGetSegmentationById_Success(t *testing.T) {
	h := newHandler(&mockSegmentationService{
		seg: domain.Segmentation{
			Id:                  7,
			SegmentationGroupID: 1,
			Points: []domain.SegmentationPoint{
				{Id: 1, SegmentationID: 7, X: 5, Y: 6},
			},
			CreateAt: time.Now().UTC(),
		},
	})

	resp, err := h.GetSegmentationById(context.Background(), &pb.GetSegmentationByIdIn{Id: 7})

	require.NoError(t, err)
	require.Equal(t, int32(7), resp.Segmentation.Id)
	require.Len(t, resp.Segmentation.Points, 1)
}

func TestGetSegmentationById_InternalError(t *testing.T) {
	h := newHandler(&mockSegmentationService{segErr: errors.New("db error")})

	_, err := h.GetSegmentationById(context.Background(), &pb.GetSegmentationByIdIn{Id: 1})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}

func TestGetSegmentsByGroupId_EmptyOnError(t *testing.T) {
	h := newHandler(&mockSegmentationService{segsErr: errors.New("db error")})

	resp, err := h.GetSegmentsByGroupId(context.Background(), &pb.GetSegmentsByGroupIdIn{SegmentationGroupId: 1})

	require.NoError(t, err)
	require.Empty(t, resp.Segmentations)
}

func TestGetSegmentsByGroupId_Success(t *testing.T) {
	h := newHandler(&mockSegmentationService{
		segs: []domain.Segmentation{
			{Id: 1, SegmentationGroupID: 2, CreateAt: time.Now().UTC()},
			{Id: 2, SegmentationGroupID: 2, CreateAt: time.Now().UTC()},
		},
	})

	resp, err := h.GetSegmentsByGroupId(context.Background(), &pb.GetSegmentsByGroupIdIn{SegmentationGroupId: 2})

	require.NoError(t, err)
	require.Len(t, resp.Segmentations, 2)
}

func TestUpdateSegmentation_Success(t *testing.T) {
	h := newHandler(&mockSegmentationService{
		updateSeg: domain.Segmentation{
			Id:                  3,
			SegmentationGroupID: 1,
			Points: []domain.SegmentationPoint{
				{X: 1, Y: 2},
			},
			CreateAt: time.Now().UTC(),
		},
	})

	resp, err := h.UpdateSegmentation(context.Background(), &pb.UpdateSegmentationIn{
		Id: 3,
		Points: []*pb.SegmentationPointCreate{
			{X: 1, Y: 2},
		},
	})

	require.NoError(t, err)
	require.Equal(t, int32(3), resp.Segmentation.Id)
}

func TestUpdateSegmentation_InternalError(t *testing.T) {
	h := newHandler(&mockSegmentationService{updateErr: errors.New("db error")})

	_, err := h.UpdateSegmentation(context.Background(), &pb.UpdateSegmentationIn{Id: 1})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}

func TestDeleteSegmentation_Success(t *testing.T) {
	h := newHandler(&mockSegmentationService{})

	resp, err := h.DeleteSegmentation(context.Background(), &pb.DeleteSegmentationIn{Id: 1})

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestDeleteSegmentation_InternalError(t *testing.T) {
	h := newHandler(&mockSegmentationService{deleteErr: errors.New("db error")})

	_, err := h.DeleteSegmentation(context.Background(), &pb.DeleteSegmentationIn{Id: 1})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
