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
		"family public.nonempty_text_field[], " +
		"upper_notes public.nonempty_text_field[], " +
		"middle_notes public.nonempty_text_field[], " +
		"base_notes public.nonempty_text_field[], " +
		"image_url TEXT, " +
		"sex_id INTEGER, " +
		"PRIMARY KEY (brand, name, sex_id), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id)" +
		");"
	CreateLinksTable = "CREATE TABLE IF NOT EXISTS perfume_links " +
		"(" +
		"brand public.nonempty_text_field, " +
		"name public.nonempty_text_field, " +
		"link public.nonempty_text_field, " +
		"volume INTEGER NOT NULL, " +
		"sex_id INTEGER, " +
		"PRIMARY KEY (brand, name, sex_id, volume), " +
		"FOREIGN KEY (brand, name, sex_id) REFERENCES perfumes(brand, name, sex_id), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id)" +
		");"
	CreateSexesTable = "CREATE TABLE IF NOT EXISTS sexes " +
		"(" +
		"sex public.nonempty_text_field UNIQUE, " +
		"id SERIAL, " +
		"PRIMARY KEY (id)" +
		");"
	FillSexesTable = "INSERT INTO sexes (sex) VALUES ('unisex'), ('female'), ('male');"

	UpdatePerfumes = "INSERT INTO perfumes (" +
		"brand, name, sex_id, perfume_type, family, upper_notes, middle_notes, base_notes, image_url" +
		") VALUES (" +
		"$1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), $4, $5, $6, $7, $8, $9" +
		") " +
		"ON CONFLICT (brand, name, sex_id) DO UPDATE SET " +
		"perfume_type = EXCLUDED.perfume_type, " +
		"family = EXCLUDED.family, " +
		"upper_notes = EXCLUDED.upper_notes, " +
		"middle_notes = EXCLUDED.middle_notes, " +
		"base_notes = EXCLUDED.base_notes, " +
		"image_url = EXCLUDED.image_url"
	UpdatePerfumeLinks = "INSERT INTO perfume_links " +
		"(" +
		"brand, name, sex_id, link, volume" +
		") VALUES (" +
		"$1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), $4, $5" +
		") " +
		"ON CONFLICT (brand, name, sex_id, volume) DO UPDATE SET " +
		"link = EXCLUDED.link"

	Select = "SELECT perfumes.brand, perfumes.name, perfumes.perfume_type, perfumes.family, " +
		"perfumes.upper_notes, perfumes.middle_notes, perfumes.base_notes, perfumes.image_url, " +
		"sexes.sex, perfume_links.volume, perfume_links.link " +
		"FROM perfumes " +
		"LEFT JOIN sexes ON perfumes.sex_id = sexes.id " +
		"LEFT JOIN perfume_links ON " +
		"perfumes.brand = perfume_links.brand AND " +
		"perfumes.name = perfume_links.name AND " +
		"perfumes.sex_id = perfume_links.sex_id"

	Truncate = "TRUNCATE perfumes, perfume_links;"

	Savepoint         = "SAVEPOINT perfume_update_"
	ReleaseSavepoint  = "RELEASE SAVEPOINT perfume_update_"
	RollbackSavepoint = "ROLLBACK TO SAVEPOINT perfume_update_"
)
