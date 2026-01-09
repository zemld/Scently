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

	InsertVariant = "INSERT INTO variants (canonized_brand, canonized_name, sex_id, shop_id, volume, price, link) " +
		"VALUES ($1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), (SELECT id FROM shops WHERE name = $4 LIMIT 1), $5, $6, $7) " +
		"ON CONFLICT (canonized_brand, canonized_name, sex_id, shop_id, volume) DO UPDATE SET " +
		"price = EXCLUDED.price, " +
		"link = EXCLUDED.link;"

	InsertFamily = "INSERT INTO families (canonized_brand, canonized_name, sex_id, family) " +
		"VALUES ($1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), $4) " +
		"ON CONFLICT (canonized_brand, canonized_name, sex_id, family) DO NOTHING;"

	InsertUpperNote = "INSERT INTO upper_notes (canonized_brand, canonized_name, sex_id, note) " +
		"VALUES ($1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), $4) " +
		"ON CONFLICT (canonized_brand, canonized_name, sex_id, note) DO NOTHING;"

	InsertCoreNote = "INSERT INTO core_notes (canonized_brand, canonized_name, sex_id, note) " +
		"VALUES ($1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), $4) " +
		"ON CONFLICT (canonized_brand, canonized_name, sex_id, note) DO NOTHING;"

	InsertBaseNote = "INSERT INTO base_notes (canonized_brand, canonized_name, sex_id, note) " +
		"VALUES ($1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), $4) " +
		"ON CONFLICT (canonized_brand, canonized_name, sex_id, note) DO NOTHING;"

	InsertPerfumeBaseInfo = "INSERT INTO perfume_base_info (canonized_brand, canonized_name, sex_id, brand, name, type, image_url, updated_at) " +
		"VALUES ($1, $2, (SELECT id FROM sexes WHERE sex = $3 LIMIT 1), $4, $5, $6, $7, CURRENT_TIMESTAMP) " +
		"ON CONFLICT (canonized_brand, canonized_name, sex_id) DO UPDATE SET " +
		"brand = EXCLUDED.brand, " +
		"name = EXCLUDED.name, " +
		"type = EXCLUDED.type, " +
		"image_url = EXCLUDED.image_url, " +
		"updated_at = CURRENT_TIMESTAMP;"

	DeleteOldPerfumes = `
	CREATE TEMP TABLE IF NOT EXISTS old_perfumes_temp (
		canonized_brand public.nonempty_text_field,
		canonized_name public.nonempty_text_field,
		sex_id INTEGER,
		PRIMARY KEY (canonized_brand, canonized_name, sex_id)
	) ON COMMIT DROP;
	
	INSERT INTO old_perfumes_temp (canonized_brand, canonized_name, sex_id)
	SELECT canonized_brand, canonized_name, sex_id 
	FROM perfume_base_info 
	WHERE updated_at < NOW() - INTERVAL '1 week';
	
	DELETE FROM variants 
	WHERE (canonized_brand, canonized_name, sex_id) IN (SELECT canonized_brand, canonized_name, sex_id FROM old_perfumes_temp);
	
	DELETE FROM families 
	WHERE (canonized_brand, canonized_name, sex_id) IN (SELECT canonized_brand, canonized_name, sex_id FROM old_perfumes_temp);
	
	DELETE FROM upper_notes 
	WHERE (canonized_brand, canonized_name, sex_id) IN (SELECT canonized_brand, canonized_name, sex_id FROM old_perfumes_temp);
	
	DELETE FROM core_notes 
	WHERE (canonized_brand, canonized_name, sex_id) IN (SELECT canonized_brand, canonized_name, sex_id FROM old_perfumes_temp);
	
	DELETE FROM base_notes 
	WHERE (canonized_brand, canonized_name, sex_id) IN (SELECT canonized_brand, canonized_name, sex_id FROM old_perfumes_temp);
	
	DELETE FROM perfume_base_info 
	WHERE (canonized_brand, canonized_name, sex_id) IN (SELECT canonized_brand, canonized_name, sex_id FROM old_perfumes_temp);
	
	DROP TABLE IF EXISTS old_perfumes_temp;`

	RefreshMV = `REFRESH MATERIALIZED VIEW perfume_base_info_with_pages;`
)
