-- +goose Up
-- +goose StatementBegin
CREATE MATERIALIZED VIEW IF NOT EXISTS perfume_base_info_with_pages AS
SELECT pb.canonized_brand,
    pb.canonized_name as canonized_name,
    pb.brand as brand,
    pb.name as name,
    pb.sex_id as sex_id,
    s.sex as sex,
    pb.type as type,
    pb.image_url as image_url,
    CEIL(
        ROW_NUMBER() OVER (
            ORDER BY pb.canonized_brand,
                pb.canonized_name,
                pb.sex_id
        ) / 400.0
    ) AS page_number
FROM perfume_base_info pb
    INNER JOIN sexes s ON pb.sex_id = s.id
ORDER BY pb.canonized_brand,
    pb.canonized_name,
    pb.sex_id;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_perfume_base_info_with_pages_page_number ON perfume_base_info_with_pages (
    canonized_brand,
    canonized_name,
    sex,
    page_number
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP MATERIALIZED VIEW IF EXISTS perfume_base_info_with_pages;
-- +goose StatementEnd