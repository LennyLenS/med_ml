-- +goose Up
-- +goose StatementBegin
ALTER TABLE card
    ADD COLUMN uuid uuid NOT NULL DEFAULT gen_random_uuid();

CREATE UNIQUE INDEX idx_card_uuid ON card(uuid);

COMMENT ON COLUMN card.uuid IS 'UUID карты пациента';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_card_uuid;

ALTER TABLE card
    DROP COLUMN uuid;
-- +goose StatementEnd