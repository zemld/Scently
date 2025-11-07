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
		"INNER JOIN sexes ON perfumes.sex_id = sexes.id " +
		"LEFT JOIN perfume_links ON " +
		"perfumes.brand = perfume_links.brand AND " +
		"perfumes.name = perfume_links.name AND " +
		"perfumes.sex_id = perfume_links.sex_id"

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
		"brand public.nonempty_text_field, " +
		"name public.nonempty_text_field, " +
		"sex_id INTEGER, " +
		"shop_id INTEGER, " +
		"volume INTEGER, " +
		"price INTEGER, " +
		"link TEXT, " +
		"PRIMARY KEY (brand, name, sex_id, shop_id, volume), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id), " +
		"FOREIGN KEY (shop_id) REFERENCES shops(id)" +
		");"

	CreateFamiliesTable = "CREATE TABLE IF NOT EXISTS families " +
		"(" +
		"brand public.nonempty_text_field, " +
		"name public.nonempty_text_field, " +
		"sex_id INTEGER, " +
		"family public.nonempty_text_field, " +
		"PRIMARY KEY (brand, name, sex_id, family), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id)" +
		");"

	CreateUpperNotesTable = "CREATE TABLE IF NOT EXISTS upper_notes " +
		"(" +
		"brand public.nonempty_text_field, " +
		"name public.nonempty_text_field, " +
		"sex_id INTEGER, " +
		"note public.nonempty_text_field, " +
		"PRIMARY KEY (brand, name, sex_id, note), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id)" +
		");"

	CreateCoreNotesTable = "CREATE TABLE IF NOT EXISTS core_notes " +
		"(" +
		"brand public.nonempty_text_field, " +
		"name public.nonempty_text_field, " +
		"sex_id INTEGER, " +
		"note public.nonempty_text_field, " +
		"PRIMARY KEY (brand, name, sex_id, note), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id)" +
		");"

	CreateBaseNotesTable = "CREATE TABLE IF NOT EXISTS base_notes " +
		"(" +
		"brand public.nonempty_text_field, " +
		"name public.nonempty_text_field, " +
		"sex_id INTEGER, " +
		"note public.nonempty_text_field, " +
		"PRIMARY KEY (brand, name, sex_id, note), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id)" +
		");"

	CreatePerfumeBaseInfoTable = "CREATE TABLE IF NOT EXISTS perfume_base_info " +
		"(" +
		"brand public.nonempty_text_field, " +
		"name public.nonempty_text_field, " +
		"sex_id INTEGER, " +
		"type TEXT, " +
		"image_url TEXT, " +
		"updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, " +
		"PRIMARY KEY (brand, name, sex_id), " +
		"FOREIGN KEY (sex_id) REFERENCES sexes(id)" +
		");"

	FillSexesTable = "INSERT INTO sexes (sex) VALUES ('unisex'), ('female'), ('male') " +
		"ON CONFLICT (sex) DO NOTHING;"

	GetOrInsertShop = "WITH inserted AS (" +
		"  INSERT INTO shops (name, domain) VALUES ($1, $2) " +
		"  ON CONFLICT (name, domain) DO NOTHING " +
		"  RETURNING id" +
		") " +
		"SELECT id FROM inserted " +
		"UNION ALL " +
		"SELECT id FROM shops WHERE name = $1 AND domain = $2 " +
		"LIMIT 1;"

	InsertVariant = "INSERT INTO variants (brand, name, sex_id, shop_id, volume, price, link) " +
		"VALUES ($1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), (SELECT id FROM shops WHERE name = $4 LIMIT 1), $5, $6, $7) " +
		"ON CONFLICT (brand, name, sex_id, shop_id, volume) DO UPDATE SET " +
		"price = EXCLUDED.price, " +
		"link = EXCLUDED.link;"

	InsertFamily = "INSERT INTO families (brand, name, sex_id, family) " +
		"VALUES ($1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), $4) " +
		"ON CONFLICT (brand, name, sex_id, family) DO NOTHING;"

	InsertUpperNote = "INSERT INTO upper_notes (brand, name, sex_id, note) " +
		"VALUES ($1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), $4) " +
		"ON CONFLICT (brand, name, sex_id, note) DO NOTHING;"

	InsertCoreNote = "INSERT INTO core_notes (brand, name, sex_id, note) " +
		"VALUES ($1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), $4) " +
		"ON CONFLICT (brand, name, sex_id, note) DO NOTHING;"

	InsertBaseNote = "INSERT INTO base_notes (brand, name, sex_id, note) " +
		"VALUES ($1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), $4) " +
		"ON CONFLICT (brand, name, sex_id, note) DO NOTHING;"

	InsertPerfumeBaseInfo = "INSERT INTO perfume_base_info (brand, name, sex_id, type, image_url, updated_at) " +
		"VALUES ($1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), $4, $5, CURRENT_TIMESTAMP) " +
		"ON CONFLICT (brand, name, sex_id) DO UPDATE SET " +
		"type = EXCLUDED.type, " +
		"image_url = EXCLUDED.image_url, " +
		"updated_at = CURRENT_TIMESTAMP;"

	DeleteOldPerfumes = `
		CREATE TEMP TABLE IF NOT EXISTS old_perfumes_temp (
			brand public.nonempty_text_field,
			name public.nonempty_text_field,
			sex_id INTEGER,
			PRIMARY KEY (brand, name, sex_id)
		) ON COMMIT DROP;
		
		INSERT INTO old_perfumes_temp (brand, name, sex_id)
		SELECT brand, name, sex_id 
		FROM perfume_base_info 
		WHERE updated_at < NOW() - INTERVAL '1 week';
		
		DELETE FROM variants 
		WHERE (brand, name, sex_id) IN (SELECT brand, name, sex_id FROM old_perfumes_temp);
		
		DELETE FROM families 
		WHERE (brand, name, sex_id) IN (SELECT brand, name, sex_id FROM old_perfumes_temp);
		
		DELETE FROM upper_notes 
		WHERE (brand, name, sex_id) IN (SELECT brand, name, sex_id FROM old_perfumes_temp);
		
		DELETE FROM core_notes 
		WHERE (brand, name, sex_id) IN (SELECT brand, name, sex_id FROM old_perfumes_temp);
		
		DELETE FROM base_notes 
		WHERE (brand, name, sex_id) IN (SELECT brand, name, sex_id FROM old_perfumes_temp);
		
		DELETE FROM perfume_base_info 
		WHERE (brand, name, sex_id) IN (SELECT brand, name, sex_id FROM old_perfumes_temp);
		
		DROP TABLE IF EXISTS old_perfumes_temp;
	`

	SelectUpgradedPerfume = `
		shops_with_variants AS (
			SELECT 
				v.brand,
				v.name,
				v.sex_id,
				sh.id as shop_id,
				sh.name as shop_name,
				sh.domain,
				json_agg(
					json_build_object(
						'volume', v.volume,
						'link', v.link,
						'price', v.price
					)
				) FILTER (WHERE v.volume IS NOT NULL) as variants
			FROM variants v
			INNER JOIN shops sh ON v.shop_id = sh.id
			GROUP BY v.brand, v.name, v.sex_id, sh.id, sh.name, sh.domain
		),
		aggregated_families AS (
			SELECT brand, name, sex_id, json_agg(DISTINCT family) as families
			FROM families
			GROUP BY brand, name, sex_id
		),
		aggregated_upper_notes AS (
			SELECT brand, name, sex_id, json_agg(DISTINCT note) as upper_notes
			FROM upper_notes
			GROUP BY brand, name, sex_id
		),
		aggregated_core_notes AS (
			SELECT brand, name, sex_id, json_agg(DISTINCT note) as core_notes
			FROM core_notes
			GROUP BY brand, name, sex_id
		),
		aggregated_base_notes AS (
			SELECT brand, name, sex_id, json_agg(DISTINCT note) as base_notes
			FROM base_notes
			GROUP BY brand, name, sex_id
		)
		SELECT 
			pb.brand,
			pb.name,
			s.sex,
			json_build_object(
				'type', pb.type,
				'family', COALESCE(af.families, '[]'::json),
				'upper_notes', COALESCE(aun.upper_notes, '[]'::json),
				'core_notes', COALESCE(acn.core_notes, '[]'::json),
				'base_notes', COALESCE(abn.base_notes, '[]'::json)
			) as properties,
			COALESCE(
				json_agg(
					json_build_object(
						'shop_name', swv.shop_name,
						'domain', swv.domain,
						'image_url', pb.image_url,
						'variants', swv.variants
					)
				) FILTER (WHERE swv.shop_id IS NOT NULL),
				'[]'::json
			) as shops
		FROM perfume_base_info pb
		INNER JOIN sexes s ON pb.sex_id = s.id
		LEFT JOIN shops_with_variants swv ON pb.brand = swv.brand AND pb.name = swv.name AND pb.sex_id = swv.sex_id
		LEFT JOIN aggregated_families af ON pb.brand = af.brand AND pb.name = af.name AND pb.sex_id = af.sex_id
		LEFT JOIN aggregated_upper_notes aun ON pb.brand = aun.brand AND pb.name = aun.name AND pb.sex_id = aun.sex_id
		LEFT JOIN aggregated_core_notes acn ON pb.brand = acn.brand AND pb.name = acn.name AND pb.sex_id = acn.sex_id
		LEFT JOIN aggregated_base_notes abn ON pb.brand = abn.brand AND pb.name = abn.name AND pb.sex_id = abn.sex_id
		GROUP BY pb.brand, pb.name, s.sex, pb.type, pb.image_url, af.families, aun.upper_notes, acn.core_notes, abn.base_notes
	`

	Savepoint         = "SAVEPOINT perfume_update_"
	ReleaseSavepoint  = "RELEASE SAVEPOINT perfume_update_"
	RollbackSavepoint = "ROLLBACK TO SAVEPOINT perfume_update_"
)
