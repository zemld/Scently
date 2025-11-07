package core

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/constants"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/models"
)

var (
	shopPriority = map[string]int{
		"Gold Apple": 1,
		"Randewoo":   2,
		"Letu":       3,
	}
)

func Update(ctx context.Context, params *models.UpdateParameters) models.ProcessedState {
	config := config.NewConfig()

	conn, err := pgx.Connect(ctx, config.GetConnectionString())
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return models.ProcessedState{Success: false}
	}
	defer conn.Close(ctx)

	tx, _ := conn.Begin(ctx)
	defer tx.Rollback(ctx)

	if !deleteOldPerfumes(ctx, tx) {
		log.Printf("Warning: Failed to delete old perfumes, continuing with update\n")
	}

	updateStatus := upsert(ctx, tx, params.Perfumes)

	tx.Commit(ctx)
	return updateStatus
}

func deleteOldPerfumes(ctx context.Context, tx pgx.Tx) bool {
	_, err := tx.Exec(ctx, constants.DeleteOldPerfumes)
	if err != nil {
		log.Printf("Error deleting old perfumes: %v\n", err)
		return false
	}
	log.Println("Old perfumes (older than 1 week) deleted successfully")
	return true
}

func upsert(ctx context.Context, tx pgx.Tx, perfumes []models.Perfume) models.ProcessedState {
	updateState := models.NewProcessedState()
	for i, perfume := range perfumes {
		updateSavepointStatus(ctx, tx, constants.Savepoint, i)
		if err := runUpdateQueries(ctx, tx, perfume); err != nil {
			log.Printf("Error updating perfume %s %s: %v\n", perfume.Brand, perfume.Name, err)
			updateSavepointStatus(ctx, tx, constants.RollbackSavepoint, i)
			updateState.FailedCount++
			continue
		}
		updateSavepointStatus(ctx, tx, constants.ReleaseSavepoint, i)
		updateState.SuccessfulCount++
	}

	return updateState
}

func runUpdateQueries(ctx context.Context, tx pgx.Tx, perfume models.Perfume) error {
	if err := updateShopInfo(ctx, tx, perfume); err != nil {
		return err
	}
	if err := updateFamilies(ctx, tx, perfume); err != nil {
		return err
	}
	if err := updateNotes(ctx, tx, constants.InsertUpperNote, perfume, perfume.Properties.UpperNotes); err != nil {
		return err
	}
	if err := updateNotes(ctx, tx, constants.InsertCoreNote, perfume, perfume.Properties.CoreNotes); err != nil {
		return err
	}
	if err := updateNotes(ctx, tx, constants.InsertBaseNote, perfume, perfume.Properties.BaseNotes); err != nil {
		return err
	}
	if err := updatePerfumeType(ctx, tx, perfume); err != nil {
		return err
	}
	return nil
}

func updateShopInfo(ctx context.Context, tx pgx.Tx, perfume models.Perfume) error {
	for _, shop := range perfume.Shops {
		if _, err := tx.Exec(ctx, constants.GetOrInsertShop, shop.ShopName, shop.Domain); err != nil {
			return err
		}
		for _, variant := range shop.Variants {
			if _, err := tx.Exec(ctx, constants.InsertVariant,
				perfume.Brand,
				perfume.Name,
				perfume.Sex,
				shop.ShopName,
				variant.Volume,
				variant.Price,
				variant.Link,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func updateFamilies(ctx context.Context, tx pgx.Tx, perfume models.Perfume) error {
	for _, family := range perfume.Properties.Family {
		if _, err := tx.Exec(ctx, constants.InsertFamily, perfume.Brand, perfume.Name, perfume.Sex, family); err != nil {
			return err
		}
	}
	return nil
}

func updateNotes(ctx context.Context, tx pgx.Tx, query string, perfume models.Perfume, notes []string) error {
	for _, note := range notes {
		if _, err := tx.Exec(ctx, query, perfume.Brand, perfume.Name, perfume.Sex, note); err != nil {
			return err
		}
	}
	return nil
}

func updatePerfumeType(ctx context.Context, tx pgx.Tx, perfume models.Perfume) error {
	imageUrl := getPreferredImageUrl(perfume)
	if _, err := tx.Exec(ctx, constants.InsertPerfumeBaseInfo, perfume.Brand, perfume.Name, perfume.Sex, perfume.Properties.Type, imageUrl); err != nil {
		return err
	}
	return nil
}

func getPreferredImageUrl(perfume models.Perfume) string {
	priority := 100
	imageUrl := ""
	for _, shop := range perfume.Shops {
		if priority > shopPriority[shop.ShopName] {
			priority = shopPriority[shop.ShopName]
			imageUrl = shop.ImageUrl
		}
	}
	return imageUrl
}

func updateSavepointStatus(ctx context.Context, tx pgx.Tx, cmd string, i int) {
	_, _ = tx.Exec(ctx, getSavepointQuery(cmd, i))
}

func getSavepointQuery(cmd string, i int) string {
	return fmt.Sprintf("%s%d", cmd, i)
}
