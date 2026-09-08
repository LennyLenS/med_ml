-- +goose Up
-- +goose StatementBegin
CREATE TABLE cytology_analysis
(
    id                uuid        PRIMARY KEY,
    cytology_id       uuid        NOT NULL REFERENCES cytology_image (id) ON DELETE CASCADE,
    original_image_id uuid        NOT NULL REFERENCES original_image (id) ON DELETE CASCADE,
    result             jsonb,
    created_at         timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at       timestamp
);

CREATE INDEX idx_cytology_analysis_cytology_id ON cytology_analysis(cytology_id);
CREATE INDEX idx_cytology_analysis_original_image_id ON cytology_analysis(original_image_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cytology_analysis;
-- +goose StatementEnd
