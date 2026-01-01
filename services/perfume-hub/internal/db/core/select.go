package core

import (
	"context"
	"log"

	perfumeModels "github.com/zemld/Scently/models"
	"github.com/zemld/Scently/perfume-hub/internal/errors"
	"github.com/zemld/Scently/perfume-hub/internal/models"
)

type SelectFunc func(ctx context.Context, params *models.SelectParameters) ([]perfumeModels.Perfume, models.ProcessedState)

func Select(ctx context.Context, params *models.SelectParameters) ([]perfumeModels.Perfume, models.ProcessedState) {
	rows, err := Pool.Query(ctx, params.GetQuery(), params.Unpack()...)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		return nil, models.ProcessedState{Error: errors.NewDBError("error executing query", err)}
	}
	defer rows.Close()

	processedState := models.NewProcessedState()
	var perfumes []perfumeModels.Perfume
	for rows.Next() {
		var perfume perfumeModels.Perfume
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
