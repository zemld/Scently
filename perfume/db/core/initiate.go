package core

import (
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/constants"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/internal"
)

func Initiate() {
	config := config.NewConfig()
	ctx, cancel := internal.CreateContext(config)
	defer cancel()

	conn, err := pgx.Connect(ctx, config.GetConnectionString())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, constants.CreateTable)
	if err != nil {
		log.Fatalf("Unable to create table: %v\n", err)
	}
	log.Println("Perfume table created successfully")
}
