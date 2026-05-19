package cytology_image_test

import (
	"context"
	"errors"
	"testing"

	"cytology/internal/domain"
	pb "cytology/internal/generated/grpc/service"
	"cytology/internal/server/cytology_image"
	cytologyimageservice "cytology/internal/services/cytology_image"
	"cytology/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockCytologyImageService struct {
	ids []uuid.UUID
	err error
}

func (m *mockCytologyImageService) CreateCytologyImage(context.Context, cytologyimageservice.CreateCytologyImageArg) (uuid.UUID, error) {
	panic("not implemented")
}

func (m *mockCytologyImageService) GetCytologyImageByID(context.Context, uuid.UUID) (domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyImageService) GetCytologyImagesByExternalID(context.Context, uuid.UUID) ([]domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyImageService) GetCytologyImagesByDoctorIdAndPatientId(context.Context, uuid.UUID, uuid.UUID) ([]domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyImageService) GetCytologyImagesByPatientId(context.Context, uuid.UUID) ([]domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyImageService) GetCytologyImageIdsByDoctorIdAndPatientId(context.Context, uuid.UUID, uuid.UUID) ([]uuid.UUID, error) {
	return m.ids, m.err
}

func (m *mockCytologyImageService) UpdateCytologyImage(context.Context, cytologyimageservice.UpdateCytologyImageArg) (domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyImageService) DeleteCytologyImage(context.Context, uuid.UUID) error {
	panic("not implemented")
}

func (m *mockCytologyImageService) CopyCytologyImage(context.Context, uuid.UUID) (domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyImageService) GetCytologyImageHistory(context.Context, uuid.UUID) ([]domain.CytologyImage, error) {
	panic("not implemented")
}

func newHandler(svc *mockCytologyImageService) cytology_image.CytologyImageHandler {
	return cytology_image.New(&services.Services{CytologyImage: svc})
}

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
	h := newHandler(&mockCytologyImageService{err: domain.ErrNotFound})

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
	h := newHandler(&mockCytologyImageService{err: errors.New("db error")})

	_, err := h.GetCytologyImageIdsByDoctorIdAndPatientId(context.Background(), &pb.GetCytologyImageIdsByDoctorIdAndPatientIdIn{
		DoctorId:  uuid.New().String(),
		PatientId: uuid.New().String(),
	})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
