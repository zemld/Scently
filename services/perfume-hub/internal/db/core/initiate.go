package core

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zemld/Scently/perfume-hub/internal/db/config"
	queries "github.com/zemld/Scently/perfume-hub/internal/db/query"
)

var Pool *pgxpool.Pool

func Initiate() {
	config := config.NewConfig()
	ctx, cancel := CreateContext(config)
	defer cancel()

	var err error
	Pool, err = pgxpool.New(context.Background(), config.GetConnectionString())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	if err := Pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	setupDatabase(ctx, Pool,
		queries.NonEmptyTextField,
		queries.CreateSexesTable,
		queries.FillSexesTable,
		queries.CreateShopsTable,
		queries.CreateVariantsTable,
		queries.CreateFamiliesTable,
		queries.CreateUpperNotesTable,
		queries.CreateCoreNotesTable,
		queries.CreateBaseNotesTable,
		queries.CreatePerfumeBaseInfoTable,
	)
	log.Println("Perfume table created successfully")
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}

func setupDatabase(ctx context.Context, db *pgxpool.Pool, setupQueries ...string) {
	for _, query := range setupQueries {
		if _, err := db.Exec(ctx, query); err != nil {
			log.Fatalf("Unable to execute query: %v\n", err)
		}
	}
}
