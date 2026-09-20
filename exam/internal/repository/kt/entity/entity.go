package entity

import (
	"database/sql"
	"time"

	"exam/internal/domain"

	"github.com/WantBeASleep/med_ml_lib/gtc"
	"github.com/google/uuid"
)

type KT struct {
	Id               uuid.UUID          `db:"id"`
	Checked          bool               `db:"checked"`
	Author           uuid.UUID          `db:"author"`
	DeviceID         int                `db:"device_id"`
	Status           string             `db:"status"`
	Description      sql.NullString     `db:"description"`
	CreateAt         time.Time          `db:"create_at"`
	PredictedClasses map[string]float64 `db:"predicted_classes"`
}

func (KT) FromDomain(d domain.KT) KT {
	return KT{
		Id:               d.Id,
		Checked:          d.Checked,
		Author:           d.Author,
		DeviceID:         d.DeviceID,
		Status:           d.Status.String(),
		Description:      gtc.String.PointerToSql(d.Description),
		CreateAt:         d.CreateAt,
		PredictedClasses: d.PredictedClasses,
	}
}

func (d KT) ToDomain() domain.KT {
	// TODO: обработать ошибку
	// но нигде встретиться не должна
	status, _ := domain.ExamStatus.Parse("", d.Status)

	return domain.KT{
		Id:               d.Id,
		Checked:          d.Checked,
		Author:           d.Author,
		DeviceID:         d.DeviceID,
		Status:           status,
		Description:      gtc.String.SqlToPointer(d.Description),
		CreateAt:         d.CreateAt,
		PredictedClasses: d.PredictedClasses,
	}
}

func (KT) SliceToDomain(kts []KT) []domain.KT {
	domainKTs := make([]domain.KT, 0, len(kts))
	for _, v := range kts {
		domainKTs = append(domainKTs, v.ToDomain())
	}
	return domainKTs
}
