package core

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/constants"
)

func Initiate() {
	config := config.NewConfig()
	ctx, cancel := CreateContext(config)
	defer cancel()

	conn, err := pgx.Connect(ctx, config.GetConnectionString())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	setupDatabase(ctx, conn,
		constants.NonEmptyTextField,
		constants.CreateSexesTable,
		constants.FillSexesTable,
		constants.CreateShopsTable,
		constants.CreateVariantsTable,
		constants.CreateFamiliesTable,
		constants.CreateUpperNotesTable,
		constants.CreateCoreNotesTable,
		constants.CreateBaseNotesTable,
		constants.CreatePerfumeBaseInfoTable,
	)
	log.Println("Perfume table created successfully")
}

func setupDatabase(ctx context.Context, conn *pgx.Conn, setupQueries ...string) {
	for _, query := range setupQueries {
		if _, err := conn.Exec(ctx, query); err != nil {
			log.Fatalf("Unable to execute query: %v\n", err)
		}
	}
}
