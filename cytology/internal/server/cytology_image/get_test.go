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

func TestGetCytologyImageById_InvalidID(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	_, err := h.GetCytologyImageById(context.Background(), &pb.GetCytologyImageByIdIn{Id: "invalid"})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetCytologyImageById_NotFound(t *testing.T) {
	h := newHandler(&mockCytologyImageService{imageErr: domain.ErrNotFound})

	_, err := h.GetCytologyImageById(context.Background(), &pb.GetCytologyImageByIdIn{Id: uuid.New().String()})

	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestGetCytologyImageById_Success(t *testing.T) {
	imageID := uuid.New()
	originalID := uuid.New()
	image := domain.CytologyImage{
		Id:               imageID,
		ExternalID:       uuid.New(),
		DoctorID:         uuid.New(),
		PatientID:        uuid.New(),
		DiagnosticNumber: 1,
		DiagnosDate:      time.Now().UTC(),
		IsLast:           true,
		CreateAt:         time.Now().UTC(),
	}

	h := newHandler(
		&mockCytologyImageService{image: image},
		&mockOriginalImageService{
			images: []domain.OriginalImage{{Id: originalID, CytologyID: imageID, ImagePath: "path/to/image"}},
		},
	)

	resp, err := h.GetCytologyImageById(context.Background(), &pb.GetCytologyImageByIdIn{Id: imageID.String()})

	require.NoError(t, err)
	require.Equal(t, imageID.String(), resp.CytologyImage.Id)
	require.NotNil(t, resp.OriginalImage)
	require.Equal(t, originalID.String(), resp.OriginalImage.Id)
}

func TestGetCytologyImageById_InternalError(t *testing.T) {
	h := newHandler(&mockCytologyImageService{imageErr: errors.New("db error")})

	_, err := h.GetCytologyImageById(context.Background(), &pb.GetCytologyImageByIdIn{Id: uuid.New().String()})

	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}

func TestGetCytologyImagesByExternalId_InvalidExternalID(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	_, err := h.GetCytologyImagesByExternalId(context.Background(), &pb.GetCytologyImagesByExternalIdIn{
		ExternalId: "invalid",
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetCytologyImagesByExternalId_NotFound(t *testing.T) {
	h := newHandler(&mockCytologyImageService{imagesErr: domain.ErrNotFound})

	resp, err := h.GetCytologyImagesByExternalId(context.Background(), &pb.GetCytologyImagesByExternalIdIn{
		ExternalId: uuid.New().String(),
	})

	require.NoError(t, err)
	require.Empty(t, resp.CytologyImages)
}

func TestGetCytologyImagesByExternalId_Success(t *testing.T) {
	imageID := uuid.New()
	h := newHandler(&mockCytologyImageService{
		images: []domain.CytologyImage{
			{
				Id:               imageID,
				ExternalID:       uuid.New(),
				DoctorID:         uuid.New(),
				PatientID:        uuid.New(),
				DiagnosticNumber: 1,
				DiagnosDate:      time.Now().UTC(),
				IsLast:           true,
				CreateAt:         time.Now().UTC(),
			},
		},
	})

	resp, err := h.GetCytologyImagesByExternalId(context.Background(), &pb.GetCytologyImagesByExternalIdIn{
		ExternalId: uuid.New().String(),
	})

	require.NoError(t, err)
	require.Len(t, resp.CytologyImages, 1)
	require.Equal(t, imageID.String(), resp.CytologyImages[0].Id)
}

func TestGetCytologyImagesByDoctorIdAndPatientId_InvalidDoctorID(t *testing.T) {
	h := newHandler(&mockCytologyImageService{})

	_, err := h.GetCytologyImagesByDoctorIdAndPatientId(context.Background(), &pb.GetCytologyImagesByDoctorIdAndPatientIdIn{
		DoctorId:  "invalid",
		PatientId: uuid.New().String(),
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetCytologyImagesByDoctorIdAndPatientId_NotFound(t *testing.T) {
	h := newHandler(&mockCytologyImageService{imagesErr: domain.ErrNotFound})

	resp, err := h.GetCytologyImagesByDoctorIdAndPatientId(context.Background(), &pb.GetCytologyImagesByDoctorIdAndPatientIdIn{
		DoctorId:  uuid.New().String(),
		PatientId: uuid.New().String(),
	})

	require.NoError(t, err)
	require.Empty(t, resp.CytologyImages)
}

func TestGetCytologyImagesByDoctorIdAndPatientId_Success(t *testing.T) {
	imageID := uuid.New()
	h := newHandler(&mockCytologyImageService{
		images: []domain.CytologyImage{
			{
				Id:               imageID,
				ExternalID:       uuid.New(),
				DoctorID:         uuid.New(),
				PatientID:        uuid.New(),
				DiagnosticNumber: 2,
				DiagnosDate:      time.Now().UTC(),
				IsLast:           true,
				CreateAt:         time.Now().UTC(),
			},
		},
	})

	resp, err := h.GetCytologyImagesByDoctorIdAndPatientId(context.Background(), &pb.GetCytologyImagesByDoctorIdAndPatientIdIn{
		DoctorId:  uuid.New().String(),
		PatientId: uuid.New().String(),
	})

	require.NoError(t, err)
	require.Len(t, resp.CytologyImages, 1)
	require.Equal(t, imageID.String(), resp.CytologyImages[0].Id)
}
