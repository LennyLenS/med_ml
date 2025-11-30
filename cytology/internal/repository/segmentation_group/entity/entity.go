package entity

import (
	"encoding/json"
	"time"

	"cytology/internal/domain"

	"github.com/google/uuid"
)

type SegmentationGroup struct {
	Id         uuid.UUID       `db:"id"`
	CytologyID uuid.UUID       `db:"cytology_id"`
	SegType    string          `db:"seg_type"`
	GroupType  string          `db:"group_type"`
	IsAI       bool            `db:"is_ai"`
	Details    json.RawMessage `db:"details"`
	CreateAt   time.Time       `db:"create_at"`
}

func (SegmentationGroup) FromDomain(d domain.SegmentationGroup) SegmentationGroup {
	return SegmentationGroup{
		Id:         d.Id,
		CytologyID: d.CytologyID,
		SegType:    d.SegType.String(),
		GroupType:  d.GroupType.String(),
		IsAI:       d.IsAI,
		Details:    d.Details,
		CreateAt:   d.CreateAt,
	}
}

func (d SegmentationGroup) ToDomain() domain.SegmentationGroup {
	var segType domain.SegType
	var groupType domain.GroupType
	// Простой парсинг, можно улучшить
	segType = domain.SegType(d.SegType)
	groupType = domain.GroupType(d.GroupType)

	return domain.SegmentationGroup{
		Id:         d.Id,
		CytologyID: d.CytologyID,
		SegType:    segType,
		GroupType:  groupType,
		IsAI:       d.IsAI,
		Details:    d.Details,
		CreateAt:   d.CreateAt,
	}
}

func (SegmentationGroup) SliceToDomain(groups []SegmentationGroup) []domain.SegmentationGroup {
	domainGroups := make([]domain.SegmentationGroup, 0, len(groups))
	for _, v := range groups {
		domainGroups = append(domainGroups, v.ToDomain())
	}
	return domainGroups
}
