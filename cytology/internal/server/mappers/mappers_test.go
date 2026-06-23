package mappers_test

import (
	"testing"
	"time"

	"cytology/internal/domain"
	pb "cytology/internal/generated/grpc/service"
	"cytology/internal/server/mappers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCytologyImageToProto(t *testing.T) {
	imageID := uuid.New()
	externalID := uuid.New()
	doctorID := uuid.New()
	patientID := uuid.New()
	prevID := uuid.New()
	marking := domain.DiagnosticMarkingP11
	material := domain.MaterialTypeGS
	calcitonin := 10
	details := `{"key":"value"}`
	now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)

	img := domain.CytologyImage{
		Id:                imageID,
		ExternalID:        externalID,
		DoctorID:          doctorID,
		PatientID:         patientID,
		DiagnosticNumber:  1,
		DiagnosticMarking: &marking,
		MaterialType:      &material,
		DiagnosDate:       now,
		IsLast:            true,
		Calcitonin:        &calcitonin,
		Details:           []byte(details),
		PrevID:            &prevID,
		CreateAt:          now,
	}

	proto := mappers.CytologyImageToProto(img)

	require.Equal(t, imageID.String(), proto.Id)
	require.Equal(t, externalID.String(), proto.ExternalId)
	require.Equal(t, pb.DiagnosticMarking_DIAGNOSTIC_MARKING_P11, *proto.DiagnosticMarking)
	require.Equal(t, pb.MaterialType_MATERIAL_TYPE_GS, *proto.MaterialType)
	require.Equal(t, int32(10), *proto.Calcitonin)
	require.Equal(t, details, *proto.Details)
	require.Equal(t, prevID.String(), *proto.PrevId)
}

func TestCreateCytologyImageArgFromProto(t *testing.T) {
	externalID := uuid.New()
	doctorID := uuid.New()
	patientID := uuid.New()
	marking := pb.DiagnosticMarking_DIAGNOSTIC_MARKING_L23
	material := pb.MaterialType_MATERIAL_TYPE_BP
	calcitonin := int32(5)
	details := `{"a":1}`
	contentType := "image/png"

	arg := mappers.CreateCytologyImageArgFromProto(
		&pb.CreateCytologyImageIn{
			DiagnosticNumber:  2,
			DiagnosticMarking: &marking,
			MaterialType:      &material,
			Calcitonin:        &calcitonin,
			Details:           &details,
			File:              []byte("file"),
			ContentType:       &contentType,
		},
		externalID,
		doctorID,
		patientID,
		nil,
		nil,
	)

	require.Equal(t, externalID, arg.ExternalID)
	require.Equal(t, doctorID, arg.DoctorID)
	require.Equal(t, patientID, arg.PatientID)
	require.Equal(t, 2, arg.DiagnosticNumber)
	require.Equal(t, domain.DiagnosticMarkingL23, *arg.DiagnosticMarking)
	require.Equal(t, domain.MaterialTypeBP, *arg.MaterialType)
	require.Equal(t, 5, *arg.Calcitonin)
	require.Equal(t, []byte(details), arg.Details)
	require.Equal(t, "image/png", arg.ContentType)
}

func TestSegmentationGroupToProto(t *testing.T) {
	cytologyID := uuid.New()
	details := `{"score":0.9}`
	now := time.Now().UTC()

	group := domain.SegmentationGroup{
		Id:         3,
		CytologyID: cytologyID,
		SegType:    domain.SegTypeSTM,
		GroupType:  domain.GroupTypeCL,
		IsAI:       true,
		Details:    []byte(details),
		CreateAt:   now,
	}

	proto := mappers.SegmentationGroupToProto(group)

	require.Equal(t, int32(3), proto.Id)
	require.Equal(t, cytologyID.String(), proto.CytologyId)
	require.Equal(t, pb.SegType_SEG_TYPE_STM, proto.SegType)
	require.Equal(t, pb.GroupType_GROUP_TYPE_CL, proto.GroupType)
	require.True(t, proto.IsAi)
	require.Equal(t, details, *proto.Details)
}

func TestSegmentationToProto(t *testing.T) {
	now := time.Now().UTC()
	seg := domain.Segmentation{
		Id:                  9,
		SegmentationGroupID: 4,
		Points: []domain.SegmentationPoint{
			{Id: 1, SegmentationID: 9, X: 100, Y: 200, UID: 42},
		},
		CreateAt: now,
	}

	proto := mappers.SegmentationToProto(seg)

	require.Equal(t, int32(9), proto.Id)
	require.Equal(t, int32(4), proto.SegmentationGroupId)
	require.Len(t, proto.Points, 1)
	require.Equal(t, int32(100), proto.Points[0].X)
	require.Equal(t, int32(200), proto.Points[0].Y)
	require.Equal(t, int64(42), proto.Points[0].Uid)
}

func TestOriginalImageToProto(t *testing.T) {
	imageID := uuid.New()
	cytologyID := uuid.New()
	delay := 2.5
	now := time.Now().UTC()

	img := domain.OriginalImage{
		Id:         imageID,
		CytologyID: cytologyID,
		ImagePath:  "bucket/path",
		CreateDate: now,
		DelayTime:  &delay,
		ViewedFlag: true,
	}

	proto := mappers.OriginalImageToProto(img)

	require.Equal(t, imageID.String(), proto.Id)
	require.Equal(t, cytologyID.String(), proto.CytologyId)
	require.Equal(t, "bucket/path", proto.ImagePath)
	require.Equal(t, delay, *proto.DelayTime)
	require.True(t, proto.ViewedFlag)
}
