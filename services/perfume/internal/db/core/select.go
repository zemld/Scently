package core

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/models"
)

func Select(ctx context.Context, params *SelectParameters) ([]models.Perfume, ProcessedState) {
	config := config.NewConfig()

	conn, err := pgx.Connect(ctx, config.GetConnectionString())
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, ProcessedState{Success: false}
	}
	defer conn.Close(ctx)

	log.Printf("Executing query: %s\n", params.getQuery())
	rows, err := conn.Query(ctx, params.getQuery(), params.unpack()...)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		return nil, ProcessedState{Success: false}
	}
	defer rows.Close()

	processedState := NewProcessedState()
	var perfumes []models.Perfume
	for rows.Next() {
		var perfume models.Perfume
		err := rows.Scan(
			&perfume.Brand,
			&perfume.Name,
			&perfume.Type,
			&perfume.Family,
			&perfume.UpperNotes,
			&perfume.MiddleNotes,
			&perfume.BaseNotes,
			&perfume.ImageUrl,
			&perfume.Sex,
			&perfume.Volume,
			&perfume.Link,
		)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			processedState.FailedCount++
			continue
		}
		perfumes = append(perfumes, perfume)
		processedState.SuccessfulCount++
	}
	return perfumes, processedState
}
