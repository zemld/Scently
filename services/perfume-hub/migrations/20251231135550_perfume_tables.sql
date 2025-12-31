-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sexes 
    (
		sex public.nonempty_text_field UNIQUE, 
		id SERIAL, 
		PRIMARY KEY (id)
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shops
    (
		id SERIAL,
		name public.nonempty_text_field,
		domain public.nonempty_text_field,
		PRIMARY KEY (id),
		UNIQUE (name, domain)
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS variants
    (
		canonized_brand public.nonempty_text_field,
		canonized_name public.nonempty_text_field,
		sex_id INTEGER,
		shop_id INTEGER,
		volume INTEGER,
		price INTEGER,
		link TEXT,
		PRIMARY KEY (canonized_brand, canonized_name, sex_id, shop_id, volume),
		FOREIGN KEY (sex_id) REFERENCES sexes(id),
		FOREIGN KEY (shop_id) REFERENCES shops(id)
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS families
    (
		canonized_brand public.nonempty_text_field,
		canonized_name public.nonempty_text_field,
		sex_id INTEGER,
		family public.nonempty_text_field,
		PRIMARY KEY (canonized_brand, canonized_name, sex_id, family),
		FOREIGN KEY (sex_id) REFERENCES sexes(id)
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS upper_notes
    (
		canonized_brand public.nonempty_text_field,
		canonized_name public.nonempty_text_field,
		sex_id INTEGER,
		note public.nonempty_text_field,
		PRIMARY KEY (canonized_brand, canonized_name, sex_id, note),
		FOREIGN KEY (sex_id) REFERENCES sexes(id)
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS core_notes
    (
		canonized_brand public.nonempty_text_field,
		canonized_name public.nonempty_text_field,
		sex_id INTEGER,
		note public.nonempty_text_field,
		PRIMARY KEY (canonized_brand, canonized_name, sex_id, note),
		FOREIGN KEY (sex_id) REFERENCES sexes(id)
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS base_notes
    (
		canonized_brand public.nonempty_text_field,
		canonized_name public.nonempty_text_field,
		sex_id INTEGER,
		note public.nonempty_text_field,
		PRIMARY KEY (canonized_brand, canonized_name, sex_id, note),
		FOREIGN KEY (sex_id) REFERENCES sexes(id)
    );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS perfume_base_info
    (
		canonized_brand public.nonempty_text_field,
		canonized_name public.nonempty_text_field,
		sex_id INTEGER,
		brand public.nonempty_text_field,
		name public.nonempty_text_field,
		type TEXT,
		image_url TEXT,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (canonized_brand, canonized_name, sex_id),
		FOREIGN KEY (sex_id) REFERENCES sexes(id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sexes CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS shops CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS variants CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS families CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS upper_notes CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS core_notes CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS base_notes CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS perfume_base_info CASCADE;
-- +goose StatementEnd