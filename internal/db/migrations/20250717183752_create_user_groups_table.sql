-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_groups (
    id          SERIAL        PRIMARY KEY,
    description TEXT          NOT NULL,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ            DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_groups;
-- +goose StatementEnd
