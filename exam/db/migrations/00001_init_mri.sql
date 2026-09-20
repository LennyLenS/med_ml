-- +goose Up
-- +goose StatementBegin
CREATE TABLE device
(
    id   SERIAL       PRIMARY KEY,
    name varchar(255) NOT NULL
);

COMMENT ON TABLE device IS 'Хранилище аппаратов на которых делались снимки';
COMMENT ON COLUMN device.name IS 'Название аппарата';

CREATE TABLE mri
(
    id              uuid            PRIMARY KEY,
    projection      varchar(255)    NOT NULL,
    checked         boolean         NOT NULL,
    external_id     uuid            NOT NULL,
    author          uuid            NOT NULL,
    device_id       integer         NOT NULL REFERENCES device (id),
    "status"        varchar(255)    NOT NULL,
    "description"   text,
    create_at       date            NOT NULL
);

COMMENT ON TABLE mri IS 'Хранилище описаний и характеристик мрт';
COMMENT ON COLUMN mri.projection IS 'Проекция в которой было сделано мрт';
COMMENT ON COLUMN mri.external_id IS 'Внешний идентификатор мрт';
COMMENT ON COLUMN mri.author IS 'ID автора мрт';
COMMENT ON COLUMN mri.device_id IS 'Идентификатор мрт аппарата на котором снято мрт';
COMMENT ON COLUMN mri.description IS 'Описание мрт от автора';
COMMENT ON COLUMN mri."status" IS 'Статус мрт';

CREATE TABLE image
(
    id     uuid PRIMARY KEY,
    mri_id uuid    NOT NULL REFERENCES mri (id) ON DELETE CASCADE,
    page   integer NOT NULL
);

COMMENT ON TABLE image IS 'Хранилище кадров в мрт';

CREATE TABLE node
(
    id              uuid        PRIMARY KEY,
    mri_id          uuid        NOT NULL REFERENCES mri (id) ON DELETE CASCADE,
    ai              boolean     NOT NULL,
    "validation"    varchar(255),
    knosp_012       real        NOT NULL,
    knosp_3        real        NOT NULL,
    knosp_4        real        NOT NULL,
    "description"   text
);

COMMENT ON TABLE node IS 'Хранилище узлов в мрт';
COMMENT ON COLUMN node.ai IS 'Автор узла(нейронка ли)';
COMMENT ON COLUMN node.mri_id IS 'Идентификатор мрт';
COMMENT ON COLUMN node."validation" IS 'валидация узла специалистом (null, invalid, valid). Доступно только для нейроночных узлов';
COMMENT ON COLUMN node.knosp_012 IS 'процент отношения к классу knosp_012';
COMMENT ON COLUMN node.knosp_3 IS 'процент отношения к классу knosp_3';
COMMENT ON COLUMN node.knosp_4 IS 'процент отношения к классу knosp_4';

CREATE TABLE segment
(
    id        uuid      PRIMARY KEY,
    node_id   uuid      NOT NULL REFERENCES node (id) ON DELETE CASCADE,
    image_id  uuid      NOT NULL REFERENCES image (id) ON DELETE CASCADE,
    contor    jsonb     NOT NULL,
    ai        boolean   NOT NULL,
    knosp_012 real      NOT NULL,
    knosp_3  real      NOT NULL,
    knosp_4  real      NOT NULL
);

COMMENT ON TABLE segment IS 'Хранилище сегментов в мрт';
COMMENT ON COLUMN segment.contor IS 'контур узла (JSON)';
COMMENT ON COLUMN segment.knosp_012 IS 'процент отношения к классу knosp_012';
COMMENT ON COLUMN segment.knosp_3 IS 'процент отношения к классу knosp_3';
COMMENT ON COLUMN segment.knosp_4 IS 'процент отношения к классу knosp_4';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS node CASCADE;
DROP TABLE IF EXISTS segment CASCADE;
DROP TABLE IF EXISTS image CASCADE;
DROP TABLE IF EXISTS device CASCADE;
DROP TABLE IF EXISTS mri CASCADE;
-- +goose StatementEnd


