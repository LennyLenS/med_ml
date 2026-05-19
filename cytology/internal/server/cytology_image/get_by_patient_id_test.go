package cytology_image_test

import (
	"context"
	"errors"
	"testing"
	"time"

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

type mockCytologyImageByPatientService struct {
	images []domain.CytologyImage
	err    error
}

func (m *mockCytologyImageByPatientService) CreateCytologyImage(context.Context, cytologyimageservice.CreateCytologyImageArg) (uuid.UUID, error) {
	panic("not implemented")
}

func (m *mockCytologyImageByPatientService) GetCytologyImageByID(context.Context, uuid.UUID) (domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyImageByPatientService) GetCytologyImagesByExternalID(context.Context, uuid.UUID) ([]domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyImageByPatientService) GetCytologyImagesByDoctorIdAndPatientId(context.Context, uuid.UUID, uuid.UUID) ([]domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyImageByPatientService) GetCytologyImagesByPatientId(context.Context, uuid.UUID) ([]domain.CytologyImage, error) {
	return m.images, m.err
}

func (m *mockCytologyImageByPatientService) GetCytologyImageIdsByDoctorIdAndPatientId(context.Context, uuid.UUID, uuid.UUID) ([]uuid.UUID, error) {
	panic("not implemented")
}

func (m *mockCytologyImageByPatientService) UpdateCytologyImage(context.Context, cytologyimageservice.UpdateCytologyImageArg) (domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyImageByPatientService) DeleteCytologyImage(context.Context, uuid.UUID) error {
	panic("not implemented")
}

func (m *mockCytologyImageByPatientService) CopyCytologyImage(context.Context, uuid.UUID) (domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyImageByPatientService) GetCytologyImageHistory(context.Context, uuid.UUID) ([]domain.CytologyImage, error) {
	panic("not implemented")
}

func newHandlerByPatient(svc cytologyimageservice.Service) cytology_image.CytologyImageHandler {
	return cytology_image.New(&services.Services{CytologyImage: svc})
}

func TestGetCytologyImagesByPatientId_InvalidPatientID(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	_, err := h.GetCytologyImagesByPatientId(context.Background(), &pb.GetCytologyImagesByPatientIdIn{
		PatientId: "invalid",
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetCytologyImagesByPatientId_Empty(t *testing.T) {
	h := newHandlerByPatient(&mockCytologyImageByPatientService{images: nil})

	resp, err := h.GetCytologyImagesByPatientId(context.Background(), &pb.GetCytologyImagesByPatientIdIn{
		PatientId: uuid.New().String(),
	})

	require.NoError(t, err)
	require.Empty(t, resp.CytologyImages)
}

func TestGetCytologyImagesByPatientId_Success(t *testing.T) {
	imageID := uuid.New()
	patientID := uuid.New()
	marking := domain.DiagnosticMarkingP11
	images := []domain.CytologyImage{
		{
			Id:                imageID,
			PatientID:         patientID,
			DiagnosticNumber:  1,
			DiagnosticMarking: &marking,
			DiagnosDate:       time.Now().UTC(),
			IsLast:            true,
			CreateAt:          time.Now().UTC(),
		},
	}

	h := newHandlerByPatient(&mockCytologyImageByPatientService{images: images})

	resp, err := h.GetCytologyImagesByPatientId(context.Background(), &pb.GetCytologyImagesByPatientIdIn{
		PatientId: patientID.String(),
	})

	require.NoError(t, err)
	require.Len(t, resp.CytologyImages, 1)
	require.Equal(t, imageID.String(), resp.CytologyImages[0].Id)
}

func TestGetCytologyImagesByPatientId_InternalError(t *testing.T) {
	h := newHandlerByPatient(&mockCytologyImageByPatientService{err: errors.New("db error")})

	_, err := h.GetCytologyImagesByPatientId(context.Background(), &pb.GetCytologyImagesByPatientIdIn{
		PatientId: uuid.New().String(),
	})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
