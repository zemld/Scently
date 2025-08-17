package core

import (
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/internal"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/models"
)

func Select(params *SelectParameters) []models.GluedPerfume {
	config := config.NewConfig()
	ctx, cancel := internal.CreateContext(config)
	defer cancel()

	conn, err := pgx.Connect(ctx, config.GetConnectionString())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	log.Printf("Executing query: %s\n", params.getQuery())
	rows, err := conn.Query(ctx, params.getQuery(), params.unpack()...)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
	}
	defer rows.Close()

	var perfumes []models.GluedPerfume
	for rows.Next() {
		var perfume models.GluedPerfume
		err := rows.Scan(
			&perfume.Brand,
			&perfume.Name,
			&perfume.Type,
			&perfume.Sex,
			&perfume.Family,
			&perfume.UpperNotes,
			&perfume.MiddleNotes,
			&perfume.BaseNotes,
			&perfume.Volumes,
			&perfume.Links)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			continue
		}
		perfumes = append(perfumes, perfume)
	}
	return perfumes
}
