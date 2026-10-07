package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type CytologyAnalysis struct {
	ID              uuid.UUID      `db:"id"`
	CytologyID      uuid.UUID      `db:"cytology_id"`
	OriginalImageID uuid.UUID      `db:"original_image_id"`
	Result          sql.NullString `db:"result"`
	CreatedAt       time.Time      `db:"created_at"`
	CompletedAt     sql.NullTime   `db:"completed_at"`
}
