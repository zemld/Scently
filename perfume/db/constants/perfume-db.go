package constants

const (
	NonEmptyTextField = "DO $$ BEGIN " +
		"IF NOT EXISTS (" +
		"  SELECT 1 FROM pg_type t JOIN pg_namespace n ON n.oid = t.typnamespace " +
		"  WHERE t.typname = 'nonempty_text_field' AND n.nspname = 'public'" +
		") THEN " +
		"  CREATE DOMAIN public.nonempty_text_field AS TEXT " +
		"  CHECK (VALUE IS NOT NULL AND LENGTH(btrim(VALUE)) > 0); " +
		"END IF; " +
		"END $$;"
	CreateTable = "CREATE TABLE IF NOT EXISTS perfumes " +
		"(" +
		"brand public.nonempty_text_field, " +
		"name public.nonempty_text_field, " +
		"perfume_type public.nonempty_text_field, " +
		"sex public.nonempty_text_field, " +
		"family public.nonempty_text_field, " +
		"upper_notes TEXT[] NOT NULL, " +
		"middle_notes TEXT[] NOT NULL, " +
		"base_notes TEXT[] NOT NULL, " +
		"volumes INT[] NOT NULL, " +
		"links TEXT[] NOT NULL, " +
		"PRIMARY KEY (brand, name)" +
		")"

	Update = "INSERT INTO perfumes (" +
		"brand, name, perfume_type, sex, family, upper_notes, middle_notes, base_notes, volumes, links" +
		") VALUES (" +
		"$1, $2, $3, $4, $5, $6, $7, $8, $9, $10" +
		") " +
		"ON CONFLICT (brand, name) DO UPDATE SET " +
		"perfume_type = EXCLUDED.perfume_type, " +
		"sex = EXCLUDED.sex, " +
		"family = EXCLUDED.family, " +
		"upper_notes = EXCLUDED.upper_notes, " +
		"middle_notes = EXCLUDED.middle_notes, " +
		"base_notes = EXCLUDED.base_notes, " +
		"volumes = EXCLUDED.volumes, " +
		"links = EXCLUDED.links"

	Select = "SELECT * FROM perfumes"

	Truncate = "TRUNCATE ONLY perfumes"

	Savepoint         = "SAVEPOINT perfume_update_"
	ReleaseSavepoint  = "RELEASE SAVEPOINT perfume_update_"
	RollbackSavepoint = "ROLLBACK TO SAVEPOINT perfume_update_"
)
