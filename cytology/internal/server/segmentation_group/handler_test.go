package segmentation_group_test

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

func TestCreateSegmentationGroup_InvalidCytologyID(t *testing.T) {
	h := newHandler(&mockSegmentationGroupService{})

	_, err := h.CreateSegmentationGroup(context.Background(), &pb.CreateSegmentationGroupIn{
		CytologyId: "invalid",
		SegType:    pb.SegType_SEG_TYPE_NIL,
		GroupType:  pb.GroupType_GROUP_TYPE_CE,
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCreateSegmentationGroup_Success(t *testing.T) {
	h := newHandler(&mockSegmentationGroupService{createID: 15})

	resp, err := h.CreateSegmentationGroup(context.Background(), &pb.CreateSegmentationGroupIn{
		CytologyId: uuid.New().String(),
		SegType:    pb.SegType_SEG_TYPE_NIR,
		GroupType:  pb.GroupType_GROUP_TYPE_CL,
		IsAi:       true,
	})

	require.NoError(t, err)
	require.Equal(t, int32(15), resp.Id)
}

func TestCreateSegmentationGroup_InternalError(t *testing.T) {
	h := newHandler(&mockSegmentationGroupService{createErr: errors.New("db error")})

	_, err := h.CreateSegmentationGroup(context.Background(), &pb.CreateSegmentationGroupIn{
		CytologyId: uuid.New().String(),
		SegType:    pb.SegType_SEG_TYPE_NIL,
		GroupType:  pb.GroupType_GROUP_TYPE_CE,
	})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}

func TestGetSegmentationGroupsByCytologyId_InvalidCytologyID(t *testing.T) {
	h := newHandler(&mockSegmentationGroupService{})

	_, err := h.GetSegmentationGroupsByCytologyId(context.Background(), &pb.GetSegmentationGroupsByCytologyIdIn{
		CytologyId: "invalid",
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetSegmentationGroupsByCytologyId_EmptyOnError(t *testing.T) {
	h := newHandler(&mockSegmentationGroupService{groupsErr: errors.New("db error")})

	resp, err := h.GetSegmentationGroupsByCytologyId(context.Background(), &pb.GetSegmentationGroupsByCytologyIdIn{
		CytologyId: uuid.New().String(),
	})

	require.NoError(t, err)
	require.Empty(t, resp.SegmentationGroups)
}

func TestGetSegmentationGroupsByCytologyId_Success(t *testing.T) {
	cytologyID := uuid.New()
	h := newHandler(&mockSegmentationGroupService{
		groups: []domain.SegmentationGroup{
			{
				Id:         1,
				CytologyID: cytologyID,
				SegType:    domain.SegTypeNIR,
				GroupType:  domain.GroupTypeCE,
				IsAI:       false,
				CreateAt:   time.Now().UTC(),
			},
		},
	})

	resp, err := h.GetSegmentationGroupsByCytologyId(context.Background(), &pb.GetSegmentationGroupsByCytologyIdIn{
		CytologyId: cytologyID.String(),
	})

	require.NoError(t, err)
	require.Len(t, resp.SegmentationGroups, 1)
	require.Equal(t, int32(1), resp.SegmentationGroups[0].Id)
	require.Equal(t, pb.SegType_SEG_TYPE_NIR, resp.SegmentationGroups[0].SegType)
}

func TestUpdateSegmentationGroup_Success(t *testing.T) {
	h := newHandler(&mockSegmentationGroupService{
		group: domain.SegmentationGroup{
			Id:         5,
			CytologyID: uuid.New(),
			SegType:    domain.SegTypeCNO,
			GroupType:  domain.GroupTypeME,
			CreateAt:   time.Now().UTC(),
		},
	})

	resp, err := h.UpdateSegmentationGroup(context.Background(), &pb.UpdateSegmentationGroupIn{Id: 5})

	require.NoError(t, err)
	require.Equal(t, int32(5), resp.SegmentationGroup.Id)
}

func TestUpdateSegmentationGroup_InternalError(t *testing.T) {
	h := newHandler(&mockSegmentationGroupService{updateErr: errors.New("db error")})

	_, err := h.UpdateSegmentationGroup(context.Background(), &pb.UpdateSegmentationGroupIn{Id: 1})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}

func TestDeleteSegmentationGroup_Success(t *testing.T) {
	h := newHandler(&mockSegmentationGroupService{})

	resp, err := h.DeleteSegmentationGroup(context.Background(), &pb.DeleteSegmentationGroupIn{Id: 1})

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestDeleteSegmentationGroup_InternalError(t *testing.T) {
	h := newHandler(&mockSegmentationGroupService{deleteErr: errors.New("db error")})

	_, err := h.DeleteSegmentationGroup(context.Background(), &pb.DeleteSegmentationGroupIn{Id: 1})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
