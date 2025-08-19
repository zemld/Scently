package core

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/constants"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/internal"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/models"
)

func Update(params *UpdateParameters, perfumes []models.Perfume) {
	config := config.NewConfig()
	ctx, cancel := internal.CreateContext(config)
	defer cancel()

	conn, err := pgx.Connect(ctx, config.GetConnectionString())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	tx, _ := conn.Begin(ctx)
	defer tx.Rollback(ctx)

	if params.IsTruncate {
		if !truncate(ctx, tx) {
			return
		}
	}

	upsert(ctx, tx, perfumes)

	tx.Commit(ctx)
	log.Println("Perfume table updated successfully")
}

func truncate(ctx context.Context, tx pgx.Tx) bool {
	_, err := tx.Exec(ctx, constants.Truncate)
	if err != nil {
		log.Printf("Error truncating perfume table: %v\n", err)
		return false
	}
	log.Println("Perfume table truncated successfully")
	return true
}

func upsert(ctx context.Context, tx pgx.Tx, perfumes []models.Perfume) {
	for i, perfume := range perfumes {
		updateSavepointStatus(ctx, tx, constants.Savepoint, i)
		_, err := tx.Exec(ctx, constants.Update, perfume.Unpack()...)
		if err != nil {
			log.Printf("Error updating perfume %s %s: %v\n", perfume.Brand, perfume.Name, err)
			updateSavepointStatus(ctx, tx, constants.RollbackSavepoint, i)
		}
		updateSavepointStatus(ctx, tx, constants.ReleaseSavepoint, i)
	}
}

func updateSavepointStatus(ctx context.Context, tx pgx.Tx, cmd string, i int) {
	_, _ = tx.Exec(ctx, getSavepointQuery(cmd, i))
}

func getSavepointQuery(cmd string, i int) string {
	return fmt.Sprintf("%s%d", cmd, i)
}
