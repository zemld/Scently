package core

import (
	"context"
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

	setupDatabase(ctx, conn)
	log.Println("Perfume table created successfully")
}

func setupDatabase(ctx context.Context, conn *pgx.Conn) {
	if _, err := conn.Exec(ctx, constants.NonEmptyTextField); err != nil {
		log.Fatalf("Unable to create nonempty_text_field domain: %v\n", err)
	}
	if _, err := conn.Exec(ctx, constants.CreatePerfumesTable); err != nil {
		log.Fatalf("Unable to create perfumes table: %v\n", err)
	}
	if _, err := conn.Exec(ctx, constants.CreateLinksTable); err != nil {
		log.Fatalf("Unable to create perfume_links table: %v\n", err)
	}
}
