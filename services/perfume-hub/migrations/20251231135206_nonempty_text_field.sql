-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type t JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'nonempty_text_field' AND n.nspname = 'public'
    ) THEN
        CREATE DOMAIN public.nonempty_text_field AS TEXT
        CHECK (VALUE IS NOT NULL AND LENGTH(btrim(VALUE)) > 0);
    END IF;
END $$;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP DOMAIN IF EXISTS public.nonempty_text_field;
-- +goose StatementEnd