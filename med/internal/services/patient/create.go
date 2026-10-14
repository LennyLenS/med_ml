package patient

import (
	"context"
	"errors"
	"fmt"

	"med/internal/domain"
	repoEntity "med/internal/repository/entity"
	"med/internal/repository/patient/entity"
	"med/internal/services/validation"
)

func (s *service) InsertPatient(ctx context.Context, patient domain.Patient) error {
	// Проверка валидности ОМС
	if !validation.ValidatePolicy(patient.Policy) {
		return fmt.Errorf("%w: неверный формат ОМС", domain.ErrBadRequest)
	}

	patientQuery := s.dao.NewPatientQuery(ctx)
	err := patientQuery.InsertPatient(entity.Patient{}.FromDomain(patient))
	if err != nil {
		var dbErr *repoEntity.DBConflictError
		if errors.As(err, &dbErr) {
			// Авторизация возвращает один ID для одной почты. Повторное добавление
			// допустимо только для того же полиса и не обновляет данные пациента.
			existing, getErr := patientQuery.GetPatientByID(patient.Id)
			if errors.Is(getErr, repoEntity.ErrNotFound) {
				return domain.ErrConflict
			}
			if getErr != nil {
				return fmt.Errorf("get existing patient: %w", getErr)
			}
			if existing.Email == patient.Email && existing.Policy == patient.Policy {
				return nil
			}
			return domain.ErrConflict
		}
		var valErr *repoEntity.DBValidationError
		if errors.As(err, &valErr) {
			return domain.ErrUnprocessableEntity
		}
		return err
	}
	return nil
}
