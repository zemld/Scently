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
	CreatePerfumesTable = "CREATE TABLE IF NOT EXISTS perfumes " +
		"(" +
		"brand public.nonempty_text_field, " +
		"name public.nonempty_text_field, " +
		"perfume_type public.nonempty_text_field, " +
		"sex public.nonempty_text_field, " +
		"family public.nonempty_text_field[], " +
		"upper_notes public.nonempty_text_field[], " +
		"middle_notes public.nonempty_text_field[], " +
		"base_notes public.nonempty_text_field[], " +
		"image_url TEXT, " +
		"PRIMARY KEY (brand, name)" +
		");"
	CreateLinksTable = "CREATE TABLE IF NOT EXISTS perfume_links " +
		"(" +
		"brand public.nonempty_text_field, " +
		"name public.nonempty_text_field, " +
		"link public.nonempty_text_field, " +
		"volume INTEGER NOT NULL, " +
		"PRIMARY KEY (brand, name, volume)" +
		");"

	UpdatePerfumes = "INSERT INTO perfumes (" +
		"brand, name, perfume_type, sex, family, upper_notes, middle_notes, base_notes, image_url" +
		") VALUES (" +
		"$1, $2, $3, $4, $5, $6, $7, $8, $9" +
		") " +
		"ON CONFLICT (brand, name) DO UPDATE SET " +
		"perfume_type = EXCLUDED.perfume_type, " +
		"sex = EXCLUDED.sex, " +
		"family = EXCLUDED.family, " +
		"upper_notes = EXCLUDED.upper_notes, " +
		"middle_notes = EXCLUDED.middle_notes, " +
		"base_notes = EXCLUDED.base_notes, " +
		"image_url = EXCLUDED.image_url"
	UpdatePerfumeLinks = "INSERT INTO perfume_links " +
		"(" +
		"brand, name, link, volume" +
		") VALUES (" +
		"$1, $2, $3, $4" +
		") " +
		"ON CONFLICT (brand, name, volume) DO UPDATE SET " +
		"link = EXCLUDED.link"

	Select = "SELECT perfumes.*, volume, link " +
		"FROM perfumes " +
		"LEFT JOIN perfume_links ON " +
		"perfumes.brand = perfume_links.brand AND perfumes.name = perfume_links.name"

	Truncate = "TRUNCATE perfumes, perfume_links;"

	Savepoint         = "SAVEPOINT perfume_update_"
	ReleaseSavepoint  = "RELEASE SAVEPOINT perfume_update_"
	RollbackSavepoint = "ROLLBACK TO SAVEPOINT perfume_update_"
)
