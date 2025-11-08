package queries

const (
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
	
	DROP TABLE IF EXISTS old_perfumes_temp;`
)
