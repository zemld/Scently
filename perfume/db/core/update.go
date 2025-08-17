package core

import (
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/constants"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/internal"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/models"
)

func Update(perfumes []models.GluedPerfume) {
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

	for _, perfume := range perfumes {
		_, err = tx.Exec(ctx, constants.UpdateQuery, perfume.Unpack()...)
		if err != nil {
			log.Printf("Error updating perfume %s %s: %v\n", perfume.Brand, perfume.Name, err)
		}
	}
	tx.Commit(ctx)
	log.Println("Perfume table updated successfully")
}
