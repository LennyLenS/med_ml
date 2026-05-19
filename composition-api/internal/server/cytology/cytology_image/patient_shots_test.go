package cytology_image_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"composition-api/internal/domain"
	auth_domain "composition-api/internal/domain/auth"
	cytology_domain "composition-api/internal/domain/cytology"
	med_domain "composition-api/internal/domain/med"
	api "composition-api/internal/generated/http/api"
	"composition-api/internal/server/cytology/cytology_image"
	"composition-api/internal/server/security"
	"composition-api/internal/services"
	cytology_srv "composition-api/internal/services/cytology"
	"composition-api/internal/services/patient"
)

type mockCytologyService struct {
	imagesByDoctorPatient []cytology_domain.CytologyImage
	imagesByPatient       []cytology_domain.CytologyImage
	originalImages        map[uuid.UUID][]cytology_domain.OriginalImage
	errDoctorPatient      error
	errPatient            error
}

func (m *mockCytologyService) CreateCytologyImage(context.Context, cytology_srv.CreateCytologyImageArg) (uuid.UUID, error) {
	panic("not implemented")
}

func (m *mockCytologyService) GetCytologyImageById(context.Context, uuid.UUID) (cytology_domain.CytologyImage, *cytology_domain.OriginalImage, error) {
	panic("not implemented")
}

func (m *mockCytologyService) GetCytologyImagesByExternalId(context.Context, uuid.UUID) ([]cytology_domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyService) GetCytologyImagesByDoctorIdAndPatientId(context.Context, uuid.UUID, uuid.UUID) ([]cytology_domain.CytologyImage, error) {
	return m.imagesByDoctorPatient, m.errDoctorPatient
}

func (m *mockCytologyService) GetCytologyImagesByPatientId(context.Context, uuid.UUID) ([]cytology_domain.CytologyImage, error) {
	return m.imagesByPatient, m.errPatient
}

func (m *mockCytologyService) UpdateCytologyImage(context.Context, cytology_srv.UpdateCytologyImageArg) (cytology_domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyService) DeleteCytologyImage(context.Context, uuid.UUID) error {
	panic("not implemented")
}

func (m *mockCytologyService) CreateOriginalImage(context.Context, cytology_srv.CreateOriginalImageArg) (uuid.UUID, error) {
	panic("not implemented")
}

func (m *mockCytologyService) GetOriginalImageById(context.Context, uuid.UUID) (cytology_domain.OriginalImage, error) {
	panic("not implemented")
}

func (m *mockCytologyService) GetOriginalImagesByCytologyId(_ context.Context, id uuid.UUID) ([]cytology_domain.OriginalImage, error) {
	if m.originalImages == nil {
		return nil, nil
	}
	return m.originalImages[id], nil
}

func (m *mockCytologyService) UpdateOriginalImage(context.Context, cytology_srv.UpdateOriginalImageArg) (cytology_domain.OriginalImage, error) {
	panic("not implemented")
}

func (m *mockCytologyService) CreateSegmentationGroup(context.Context, cytology_srv.CreateSegmentationGroupArg) (int, error) {
	panic("not implemented")
}

func (m *mockCytologyService) GetSegmentationGroupsByCytologyId(context.Context, uuid.UUID, *cytology_domain.SegType, *cytology_domain.GroupType, *bool) ([]cytology_domain.SegmentationGroup, error) {
	panic("not implemented")
}

func (m *mockCytologyService) UpdateSegmentationGroup(context.Context, cytology_srv.UpdateSegmentationGroupArg) (cytology_domain.SegmentationGroup, error) {
	panic("not implemented")
}

func (m *mockCytologyService) DeleteSegmentationGroup(context.Context, int) error {
	panic("not implemented")
}

func (m *mockCytologyService) CreateSegmentation(context.Context, cytology_srv.CreateSegmentationArg) (int, error) {
	panic("not implemented")
}

func (m *mockCytologyService) GetSegmentationById(context.Context, int) (cytology_domain.Segmentation, error) {
	panic("not implemented")
}

func (m *mockCytologyService) GetSegmentsByGroupId(context.Context, int) ([]cytology_domain.Segmentation, error) {
	panic("not implemented")
}

func (m *mockCytologyService) UpdateSegmentation(context.Context, cytology_srv.UpdateSegmentationArg) (cytology_domain.Segmentation, error) {
	panic("not implemented")
}

func (m *mockCytologyService) DeleteSegmentation(context.Context, int) error {
	panic("not implemented")
}

func (m *mockCytologyService) CopyCytologyImage(context.Context, uuid.UUID) (cytology_domain.CytologyImage, error) {
	panic("not implemented")
}

func (m *mockCytologyService) GetCytologyImageHistory(context.Context, uuid.UUID) ([]cytology_domain.CytologyImage, error) {
	panic("not implemented")
}

type mockPatientService struct {
	patient med_domain.Patient
	err     error
}

func (m *mockPatientService) CreatePatient(context.Context, patient.CreatePatientArg) (uuid.UUID, error) {
	panic("not implemented")
}

func (m *mockPatientService) GetPatient(context.Context, uuid.UUID) (med_domain.Patient, error) {
	return m.patient, m.err
}

func (m *mockPatientService) GetPatientsByDoctorID(context.Context, uuid.UUID, *bool) ([]med_domain.Patient, error) {
	panic("not implemented")
}

func (m *mockPatientService) UpdatePatient(context.Context, uuid.UUID, patient.UpdatePatientArg) (med_domain.Patient, error) {
	panic("not implemented")
}

type mockCardService struct {
	card med_domain.Card
	err  error
}

func (m *mockCardService) CreateCard(context.Context, med_domain.Card) (med_domain.Card, error) {
	panic("not implemented")
}

func (m *mockCardService) GetCard(context.Context, uuid.UUID, uuid.UUID) (med_domain.Card, error) {
	return m.card, m.err
}

func (m *mockCardService) GetCardByID(context.Context, int) (med_domain.Card, error) {
	panic("not implemented")
}

func (m *mockCardService) UpdateCard(context.Context, med_domain.Card) (med_domain.Card, error) {
	panic("not implemented")
}

func newPatientShotsHandler(cytologySvc *mockCytologyService, patientSvc *mockPatientService, cardSvc *mockCardService) cytology_image.CytologyImageHandler {
	return cytology_image.NewHandler(&services.Services{
		CytologyService: cytologySvc,
		PatientService:  patientSvc,
		CardService:     cardSvc,
	})
}

func ctxWithRole(userID uuid.UUID, role auth_domain.Role) context.Context {
	return security.WithToken(context.Background(), security.Token{
		Id:   userID,
		Role: role,
	})
}

func TestCytologyPatientShotsRead_Doctor_ReturnsOnlyOwnShots(t *testing.T) {
	doctorID := uuid.New()
	patientID := uuid.New()
	imageID := uuid.New()
	originalID := uuid.New()

	h := newPatientShotsHandler(
		&mockCytologyService{
			imagesByDoctorPatient: []cytology_domain.CytologyImage{
				{
					Id:               imageID,
					DoctorID:         doctorID,
					PatientID:        patientID,
					DiagnosticNumber: 1,
					DiagnosDate:      time.Now().UTC(),
					IsLast:           true,
				},
			},
			originalImages: map[uuid.UUID][]cytology_domain.OriginalImage{
				imageID: {{Id: originalID, CytologyID: imageID}},
			},
		},
		&mockPatientService{
			patient: med_domain.Patient{
				Id:       patientID,
				FullName: "Иванов Иван Иванович",
				Email:    "iii@medml.med",
				Policy:   "1234123412341234",
				Active:   true,
				BirthDate: time.Date(2026, 5, 6, 0, 0, 0, 0, time.UTC),
			},
		},
		&mockCardService{},
	)

	res, err := h.CytologyPatientShotsRead(ctxWithRole(doctorID, auth_domain.RoleDoctor), api.CytologyPatientShotsReadParams{
		PatientID: patientID,
	})
	require.NoError(t, err)

	ok, isOK := res.(*api.CytologyPatientShotsReadOK)
	require.True(t, isOK)
	require.Equal(t, patientID, ok.Patient.ID)
	require.Equal(t, "Иван", ok.Patient.FirstName)
	require.Len(t, ok.Shots, 1)
	require.Equal(t, imageID, ok.Shots[0].ID)
	require.True(t, ok.Shots[0].OriginalImage.Set)
	require.Equal(t, originalID, ok.Shots[0].OriginalImage.Value)
}

func TestCytologyPatientShotsRead_Patient_ReturnsAllShots(t *testing.T) {
	patientID := uuid.New()
	imageID := uuid.New()

	h := newPatientShotsHandler(
		&mockCytologyService{
			imagesByPatient: []cytology_domain.CytologyImage{
				{
					Id:               imageID,
					PatientID:        patientID,
					DiagnosticNumber: 1,
					DiagnosDate:      time.Now().UTC(),
					IsLast:           true,
				},
			},
		},
		&mockPatientService{
			patient: med_domain.Patient{
				Id:       patientID,
				FullName: "Иванов Иван Иванович",
				Active:   true,
				BirthDate: time.Date(2026, 5, 6, 0, 0, 0, 0, time.UTC),
			},
		},
		&mockCardService{},
	)

	res, err := h.CytologyPatientShotsRead(ctxWithRole(patientID, auth_domain.RolePatient), api.CytologyPatientShotsReadParams{
		PatientID: patientID,
	})
	require.NoError(t, err)

	ok, isOK := res.(*api.CytologyPatientShotsReadOK)
	require.True(t, isOK)
	require.Len(t, ok.Shots, 1)
	require.Equal(t, imageID, ok.Shots[0].ID)
}

func TestCytologyPatientShotsRead_Patient_ForbiddenForOtherPatient(t *testing.T) {
	h := newPatientShotsHandler(&mockCytologyService{}, &mockPatientService{}, &mockCardService{})

	res, err := h.CytologyPatientShotsRead(ctxWithRole(uuid.New(), auth_domain.RolePatient), api.CytologyPatientShotsReadParams{
		PatientID: uuid.New(),
	})
	require.NoError(t, err)

	forbidden, isForbidden := res.(*api.CytologyPatientShotsReadForbidden)
	require.True(t, isForbidden)
	require.Equal(t, http.StatusForbidden, forbidden.StatusCode)
}

func TestCytologyPatientShotsRead_PatientNotFound(t *testing.T) {
	doctorID := uuid.New()
	patientID := uuid.New()

	h := newPatientShotsHandler(
		&mockCytologyService{},
		&mockPatientService{err: domain.ErrNotFound},
		&mockCardService{},
	)

	res, err := h.CytologyPatientShotsRead(ctxWithRole(doctorID, auth_domain.RoleDoctor), api.CytologyPatientShotsReadParams{
		PatientID: patientID,
	})
	require.NoError(t, err)

	notFound, isNotFound := res.(*api.CytologyPatientShotsReadNotFound)
	require.True(t, isNotFound)
	require.Equal(t, http.StatusNotFound, notFound.StatusCode)
}

func TestCytologyPatientShotsRead_UnauthorizedWithoutToken(t *testing.T) {
	h := newPatientShotsHandler(&mockCytologyService{}, &mockPatientService{}, &mockCardService{})

	_, err := h.CytologyPatientShotsRead(context.Background(), api.CytologyPatientShotsReadParams{
		PatientID: uuid.New(),
	})
	require.ErrorIs(t, err, security.ErrUnauthorized)
}
