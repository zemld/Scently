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
)
