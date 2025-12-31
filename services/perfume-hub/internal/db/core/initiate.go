package core

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/zemld/Scently/perfume-hub/internal/db/config"
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

	if err := runMigrations(ctx, config.GetConnectionString()); err != nil {
		log.Fatalf("Unable to run migrations: %v\n", err)
	}
	log.Println("Migrations run successfully")
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}

func runMigrations(ctx context.Context, connectionString string) error {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	return goose.UpContext(ctx, db, "./migrations")
}
