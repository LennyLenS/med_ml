package original_image_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	dbusproducers "cytology/internal/dbus/producers"
	"cytology/internal/domain"
	analysisrequestedpb "cytology/internal/generated/dbus/produce/cytologyanalysisrequested"
	grpcpb "cytology/internal/generated/grpc/service"
	"cytology/internal/server/original_image"
	"cytology/internal/services"
	"cytology/internal/services/cytology_analysis"
	"cytology/internal/services/cytology_image"
)

type producerSpy struct {
	message *analysisrequestedpb.AnalysisRequested
}

func (p *producerSpy) SendCytologyAnalysisRequested(_ context.Context, message *analysisrequestedpb.AnalysisRequested) error {
	p.message = message
	return nil
}

var _ dbusproducers.Producer = (*producerSpy)(nil)

type analysisServiceStub struct {
	id uuid.UUID
}

func (s analysisServiceStub) Create(context.Context, uuid.UUID, uuid.UUID) (uuid.UUID, error) {
	return s.id, nil
}

func (analysisServiceStub) Get(context.Context, uuid.UUID) (cytology_analysis.Analysis, error) {
	return cytology_analysis.Analysis{}, nil
}

func (analysisServiceStub) Complete(context.Context, uuid.UUID, []byte) error {
	return nil
}

type cytologyImageServiceStub struct {
	image domain.CytologyImage
}

func (s cytologyImageServiceStub) GetCytologyImageByID(context.Context, uuid.UUID) (domain.CytologyImage, error) {
	return s.image, nil
}

func (cytologyImageServiceStub) CreateCytologyImage(context.Context, cytology_image.CreateCytologyImageArg) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (cytologyImageServiceStub) GetCytologyImagesByExternalID(context.Context, uuid.UUID) ([]domain.CytologyImage, error) {
	return nil, nil
}

func (cytologyImageServiceStub) GetCytologyImagesByDoctorIdAndPatientId(context.Context, uuid.UUID, uuid.UUID) ([]domain.CytologyImage, error) {
	return nil, nil
}

func (cytologyImageServiceStub) GetCytologyImagesByPatientId(context.Context, uuid.UUID) ([]domain.CytologyImage, error) {
	return nil, nil
}

func (cytologyImageServiceStub) GetCytologyImageIdsByDoctorIdAndPatientId(context.Context, uuid.UUID, uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}

func (cytologyImageServiceStub) UpdateCytologyImage(context.Context, cytology_image.UpdateCytologyImageArg) (domain.CytologyImage, error) {
	return domain.CytologyImage{}, nil
}

func (cytologyImageServiceStub) DeleteCytologyImage(context.Context, uuid.UUID) error {
	return nil
}

func (cytologyImageServiceStub) CopyCytologyImage(context.Context, uuid.UUID) (domain.CytologyImage, error) {
	return domain.CytologyImage{}, nil
}

func (cytologyImageServiceStub) GetCytologyImageHistory(context.Context, uuid.UUID) ([]domain.CytologyImage, error) {
	return nil, nil
}

func TestCreateOriginalImage_SendsAnalysisRequestedForLNP(t *testing.T) {
	cytologyID := uuid.New()
	originalImageID := uuid.New()
	analysisID := uuid.New()
	materialType := domain.MaterialTypeLNP
	calcitonin := 17
	thyroglobulin := 29
	imagePath := cytologyID.String() + "/" + originalImageID.String() + "/slide.svs"
	producer := &producerSpy{}

	h := original_image.New(&services.Services{
		OriginalImage: &mockOriginalImageService{
			createID: originalImageID,
			image:    domain.OriginalImage{Id: originalImageID, CytologyID: cytologyID, ImagePath: imagePath},
		},
		CytologyImage: cytologyImageServiceStub{image: domain.CytologyImage{
			Id: cytologyID, MaterialType: &materialType, CalcitoninInFlush: &calcitonin, Thyroglobulin: &thyroglobulin,
		}},
		CytologyAnalysis: analysisServiceStub{id: analysisID},
		Producer:         producer,
	})

	response, err := h.CreateOriginalImage(context.Background(), &grpcpb.CreateOriginalImageIn{
		CytologyId: cytologyID.String(),
		ImagePath:  &imagePath,
	})

	require.NoError(t, err)
	require.Equal(t, originalImageID.String(), response.Id)
	require.NotNil(t, producer.message)
	require.Equal(t, analysisID.String(), producer.message.AnalysisId)
	require.Equal(t, "s3://cytology/"+imagePath, producer.message.InputArtifactUri)
	require.Equal(t, analysisrequestedpb.MaterialType_MATERIAL_TYPE_LNP, producer.message.MaterialType)
	require.Equal(t, int32(17), producer.message.GetCalcitoninInFlush())
	require.Equal(t, int32(29), producer.message.GetThyroglobulinInFlush())
}
