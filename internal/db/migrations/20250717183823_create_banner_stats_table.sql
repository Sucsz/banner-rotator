-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banner_stats (
    banner_id     INT         NOT NULL REFERENCES banners(id),
    slot_id       INT         NOT NULL REFERENCES slots(id),
    user_group_id INT         NOT NULL REFERENCES user_groups(id),
    impressions   BIGINT      NOT NULL DEFAULT 0,
    clicks        BIGINT      NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (banner_id, slot_id, user_group_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS banner_stats;
-- +goose StatementEnd
