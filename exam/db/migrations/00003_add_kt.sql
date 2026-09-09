-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    kt (
        id uuid PRIMARY KEY,
        checked boolean NOT NULL,
        author uuid NOT NULL,
        device_id integer NOT NULL REFERENCES device (id),
        "status" varchar(255) NOT NULL,
        "description" text,
        create_at date NOT NULL,
        predicted_classes JSONB
    );

COMMENT ON TABLE kt IS 'Хранилище описаний, характеристик и результатов анализа КТ';

COMMENT ON COLUMN kt.author IS 'ID автора КТ';

COMMENT ON COLUMN kt.device_id IS 'Идентификатор аппарата на котором снято КТ';

COMMENT ON COLUMN kt.description IS 'Описание КТ от автора';

COMMENT ON COLUMN kt."status" IS 'Статус КТ';
COMMENT ON COLUMN kt.predicted_classes IS 'Предсказанные классы КТ в виде словаря';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS kt CASCADE;
-- +goose StatementEnd+

