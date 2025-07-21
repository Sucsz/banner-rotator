-- +goose Up
-- +goose StatementBegin
INSERT INTO banners (id, title, content, description)
VALUES
    (1, 'Banner A', 'Buy now!', 'Banner for promo A'),
    (2, 'Banner B', 'Sale today!', 'Banner for promo B');

INSERT INTO slots (id, description)
VALUES
    (1, 'Main Page Slot'),
    (2, 'Sidebar Slot');

INSERT INTO user_groups (id, description)
VALUES
    (1, 'Guest users'),
    (2, 'Logged-in users');

INSERT INTO banner_slots (banner_id, slot_id)
VALUES
    (2, 1),
    (2, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM banner_slots WHERE banner_id IN (1, 2);
DELETE FROM user_groups WHERE id IN (1, 2);
DELETE FROM slots WHERE id IN (1, 2);
DELETE FROM banners WHERE id IN (1, 2);
-- +goose StatementEnd
