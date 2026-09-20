package entity

import (
	"database/sql"
	"time"

	"exam/internal/domain"

	"github.com/WantBeASleep/med_ml_lib/gtc"
	"github.com/google/uuid"
)

type Mri struct {
	Id          uuid.UUID      `db:"id"`
	Projection  string         `db:"projection"`
	Checked     bool           `db:"checked"`
	ExternalID  uuid.UUID      `db:"external_id"`
	Author      uuid.UUID      `db:"author"`
	DeviceID    int            `db:"device_id"`
	Status      string         `db:"status"`
	Description sql.NullString `db:"description"`
	CreateAt    time.Time      `db:"create_at"`
}

func (Mri) FromDomain(d domain.Mri) Mri {
	return Mri{
		Id:          d.Id,
		Projection:  d.Projection.String(),
		Checked:     d.Checked,
		ExternalID:  d.ExternalID,
		Author:      d.Author,
		DeviceID:    d.DeviceID,
		Status:      d.Status.String(),
		Description: gtc.String.PointerToSql(d.Description),
		CreateAt:    d.CreateAt,
	}
}

func (d Mri) ToDomain() domain.Mri {
	// TODO: обработать ошибку
	// но нигде встретиться не должна
	status, _ := domain.ExamStatus.Parse("", d.Status)
	projection, _ := domain.MriProjection.Parse("", d.Projection)

	return domain.Mri{
		Id:          d.Id,
		Projection:  projection,
		Checked:     d.Checked,
		ExternalID:  d.ExternalID,
		Author:      d.Author,
		DeviceID:    d.DeviceID,
		Status:      status,
		Description: gtc.String.SqlToPointer(d.Description),
		CreateAt:    d.CreateAt,
	}
}

func (Mri) SliceToDomain(mris []Mri) []domain.Mri {
	domainMris := make([]domain.Mri, 0, len(mris))
	for _, v := range mris {
		domainMris = append(domainMris, v.ToDomain())
	}
	return domainMris
}
