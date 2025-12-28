package core

import (
	"context"
	"log"

	"github.com/zemld/Scently/perfume-hub/internal/errors"
	"github.com/zemld/Scently/perfume-hub/internal/models"
)

func Select(ctx context.Context, params *models.SelectParameters) ([]models.Perfume, models.ProcessedState) {
	rows, err := Pool.Query(ctx, params.GetQuery(), params.Unpack()...)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		return nil, models.ProcessedState{Error: errors.NewDBError("error executing query", err)}
	}
	defer rows.Close()

	processedState := models.NewProcessedState()
	var perfumes []models.Perfume
	for rows.Next() {
		var perfume models.Perfume
		err := rows.Scan(
			&perfume.Brand,
			&perfume.Name,
			&perfume.Sex,
			&perfume.ImageUrl,
			&perfume.Properties,
			&perfume.Shops,
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
