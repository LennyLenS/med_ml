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

func TestGetCytologyImageIdsByDoctorIdAndPatientId_InvalidDoctorID(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	_, err := h.GetCytologyImageIdsByDoctorIdAndPatientId(context.Background(), &pb.GetCytologyImageIdsByDoctorIdAndPatientIdIn{
		DoctorId:  "invalid",
		PatientId: uuid.New().String(),
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetCytologyImageIdsByDoctorIdAndPatientId_InvalidPatientID(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	_, err := h.GetCytologyImageIdsByDoctorIdAndPatientId(context.Background(), &pb.GetCytologyImageIdsByDoctorIdAndPatientIdIn{
		DoctorId:  uuid.New().String(),
		PatientId: "invalid",
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetCytologyImageIdsByDoctorIdAndPatientId_NotFound(t *testing.T) {
	h := newHandler(&mockCytologyImageService{idsErr: domain.ErrNotFound})

	resp, err := h.GetCytologyImageIdsByDoctorIdAndPatientId(context.Background(), &pb.GetCytologyImageIdsByDoctorIdAndPatientIdIn{
		DoctorId:  uuid.New().String(),
		PatientId: uuid.New().String(),
	})

	require.NoError(t, err)
	require.Empty(t, resp.Ids)
}

func TestGetCytologyImageIdsByDoctorIdAndPatientId_Success(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	h := newHandler(&mockCytologyImageService{ids: []uuid.UUID{id1, id2}})

	resp, err := h.GetCytologyImageIdsByDoctorIdAndPatientId(context.Background(), &pb.GetCytologyImageIdsByDoctorIdAndPatientIdIn{
		DoctorId:  uuid.New().String(),
		PatientId: uuid.New().String(),
	})

	require.NoError(t, err)
	require.Equal(t, []string{id1.String(), id2.String()}, resp.Ids)
}

func TestGetCytologyImageIdsByDoctorIdAndPatientId_InternalError(t *testing.T) {
	h := newHandler(&mockCytologyImageService{idsErr: errors.New("db error")})

	_, err := h.GetCytologyImageIdsByDoctorIdAndPatientId(context.Background(), &pb.GetCytologyImageIdsByDoctorIdAndPatientIdIn{
		DoctorId:  uuid.New().String(),
		PatientId: uuid.New().String(),
	})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
