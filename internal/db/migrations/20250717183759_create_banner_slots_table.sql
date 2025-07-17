-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banner_slots (
    banner_id   INT         NOT NULL REFERENCES banners(id),
    slot_id     INT         NOT NULL REFERENCES slots(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (banner_id, slot_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS banner_slots;
-- +goose StatementEnd
