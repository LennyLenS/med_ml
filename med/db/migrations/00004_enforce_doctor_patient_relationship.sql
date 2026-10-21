-- +goose Up
-- +goose StatementBegin
-- Проверяем связи карт с врачами и пациентами без удаления данных.
-- Некорректные связи приводят к откату миграции.
ALTER TABLE card
    ADD CONSTRAINT fk_card_doctor FOREIGN KEY (doctor_id) REFERENCES doctor(id),
    ADD CONSTRAINT fk_card_patient FOREIGN KEY (patient_id) REFERENCES patient(id);

-- Добавляем индекс для поиска карт по пациенту.
CREATE INDEX idx_card_patient_doctor ON card(patient_id, doctor_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_card_patient_doctor;
ALTER TABLE card
    DROP CONSTRAINT fk_card_patient,
    DROP CONSTRAINT fk_card_doctor;
-- +goose StatementEnd
