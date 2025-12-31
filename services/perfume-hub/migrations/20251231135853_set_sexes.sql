-- +goose Up
-- +goose StatementBegin
INSERT INTO sexes (sex) VALUES ('unisex'), ('female'), ('male')
    ON CONFLICT (sex) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM sexes WHERE sex IN ('unisex', 'female', 'male');
-- +goose StatementEnd