-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banners (
    id          SERIAL        PRIMARY KEY,
    title       TEXT          NOT NULL,
    content     TEXT          NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ            DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS banners;
-- +goose StatementEnd
