package constants

const (
	CreateTable = "CREATE TABLE IF NOT EXISTS perfumes " +
		"(" +
		"brand TEXT NOT NULL, " +
		"name TEXT NOT NULL, " +
		"perfume_type TEXT NOT NULL, " +
		"sex TEXT NOT NULL, " +
		"family TEXT NOT NULL, " +
		"upper_notes TEXT[] NOT NULL, " +
		"middle_notes TEXT[] NOT NULL, " +
		"base_notes TEXT[] NOT NULL, " +
		"volumes INT[] NOT NULL, " +
		"links TEXT[] NOT NULL, " +
		"PRIMARY KEY (brand, name)" +
		")"

	UpdateQuery = "INSERT INTO perfumes (" +
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
)
