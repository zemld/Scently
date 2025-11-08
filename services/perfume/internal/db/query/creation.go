package queries

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

	CreateSexesTable = "CREATE TABLE IF NOT EXISTS sexes " +
		"(" +
		"sex public.nonempty_text_field UNIQUE, " +
		"id SERIAL, " +
		"PRIMARY KEY (id)" +
		");"

	CreateShopsTable = "CREATE TABLE IF NOT EXISTS shops " +
		"(" +
		"id SERIAL, " +
		"name public.nonempty_text_field, " +
		"domain public.nonempty_text_field, " +
		"PRIMARY KEY (id), " +
		"UNIQUE (name, domain)" +
		");"

	CreateVariantsTable = "CREATE TABLE IF NOT EXISTS variants " +
		"(" +
		"canonized_brand public.nonempty_text_field, " +
		"canonized_name public.nonempty_text_field, " +
		"sex_id INTEGER, " +
		"shop_id INTEGER, " +
		"volume INTEGER, " +
		"price INTEGER, " +
		"link TEXT, " +
		"PRIMARY KEY (canonized_brand, canonized_name, sex_id, shop_id, volume), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id), " +
		"FOREIGN KEY (shop_id) REFERENCES shops(id)" +
		");"

	CreateFamiliesTable = "CREATE TABLE IF NOT EXISTS families " +
		"(" +
		"canonized_brand public.nonempty_text_field, " +
		"canonized_name public.nonempty_text_field, " +
		"sex_id INTEGER, " +
		"family public.nonempty_text_field, " +
		"PRIMARY KEY (canonized_brand, canonized_name, sex_id, family), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id)" +
		");"

	CreateUpperNotesTable = "CREATE TABLE IF NOT EXISTS upper_notes " +
		"(" +
		"canonized_brand public.nonempty_text_field, " +
		"canonized_name public.nonempty_text_field, " +
		"sex_id INTEGER, " +
		"note public.nonempty_text_field, " +
		"PRIMARY KEY (canonized_brand, canonized_name, sex_id, note), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id)" +
		");"

	CreateCoreNotesTable = "CREATE TABLE IF NOT EXISTS core_notes " +
		"(" +
		"canonized_brand public.nonempty_text_field, " +
		"canonized_name public.nonempty_text_field, " +
		"sex_id INTEGER, " +
		"note public.nonempty_text_field, " +
		"PRIMARY KEY (canonized_brand, canonized_name, sex_id, note), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id)" +
		");"

	CreateBaseNotesTable = "CREATE TABLE IF NOT EXISTS base_notes " +
		"(" +
		"canonized_brand public.nonempty_text_field, " +
		"canonized_name public.nonempty_text_field, " +
		"sex_id INTEGER, " +
		"note public.nonempty_text_field, " +
		"PRIMARY KEY (canonized_brand, canonized_name, sex_id, note), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id)" +
		");"

	CreatePerfumeBaseInfoTable = "CREATE TABLE IF NOT EXISTS perfume_base_info " +
		"(" +
		"canonized_brand public.nonempty_text_field, " +
		"canonized_name public.nonempty_text_field, " +
		"sex_id INTEGER, " +
		"brand public.nonempty_text_field, " +
		"name public.nonempty_text_field, " +
		"type TEXT, " +
		"image_url TEXT, " +
		"updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, " +
		"PRIMARY KEY (canonized_brand, canonized_name, sex_id), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id)" +
		");"

	FillSexesTable = "INSERT INTO sexes (sex) VALUES ('unisex'), ('female'), ('male') " +
		"ON CONFLICT (sex) DO NOTHING;"
)
