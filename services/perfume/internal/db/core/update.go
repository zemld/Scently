package core

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	queries "github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/query"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/errors"
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
	tx, err := Pool.Begin(ctx)
	if err != nil {
		log.Printf("Unable to begin transaction: %v\n", err)
		return models.ProcessedState{Error: errors.NewDBError("unable to begin transaction", err)}
	}
	defer tx.Rollback(ctx)

	if !deleteOldPerfumes(ctx, tx) {
		log.Printf("Warning: Failed to delete old perfumes, continuing with update\n")
	}

	updateStatus := upsert(ctx, tx, params.Perfumes)

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Unable to commit transaction: %v\n", err)
		return models.ProcessedState{Error: errors.NewDBError("unable to commit transaction", err)}
	}
	return updateStatus
}

func deleteOldPerfumes(ctx context.Context, tx pgx.Tx) bool {
	_, err := tx.Exec(ctx, queries.DeleteOldPerfumes)
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
		updateSavepointStatus(ctx, tx, queries.Savepoint, i)
		if err := runUpdateQueries(ctx, tx, perfume); err != nil {
			log.Printf("Error updating perfume %s %s: %v\n", perfume.Brand, perfume.Name, err)
			updateSavepointStatus(ctx, tx, queries.RollbackSavepoint, i)
			updateState.FailedCount++
			continue
		}
		updateSavepointStatus(ctx, tx, queries.ReleaseSavepoint, i)
		updateState.SuccessfulCount++
	}

	return updateState
}

func runUpdateQueries(ctx context.Context, tx pgx.Tx, perfume models.Perfume) error {
	canonizedPerfume := perfume.Canonize()
	if err := updateShopInfo(ctx, tx, perfume, canonizedPerfume); err != nil {
		return err
	}
	if err := updateFamilies(ctx, tx, perfume, canonizedPerfume); err != nil {
		return err
	}
	if err := updateNotes(ctx, tx, queries.InsertUpperNote, perfume, canonizedPerfume, perfume.Properties.UpperNotes); err != nil {
		return err
	}
	if err := updateNotes(ctx, tx, queries.InsertCoreNote, perfume, canonizedPerfume, perfume.Properties.CoreNotes); err != nil {
		return err
	}
	if err := updateNotes(ctx, tx, queries.InsertBaseNote, perfume, canonizedPerfume, perfume.Properties.BaseNotes); err != nil {
		return err
	}
	if err := updatePerfumeType(ctx, tx, perfume, canonizedPerfume); err != nil {
		return err
	}
	return nil
}

func updateShopInfo(ctx context.Context, tx pgx.Tx, perfume models.Perfume, canonizedPerfume models.CanonizedPerfume) error {
	for _, shop := range perfume.Shops {
		if _, err := tx.Exec(ctx, queries.GetOrInsertShop, shop.ShopName, shop.Domain); err != nil {
			return err
		}
		for _, variant := range shop.Variants {
			if _, err := tx.Exec(ctx, queries.InsertVariant,
				canonizedPerfume.Brand,
				canonizedPerfume.Name,
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

func updateFamilies(ctx context.Context, tx pgx.Tx, perfume models.Perfume, canonizedPerfume models.CanonizedPerfume) error {
	for _, family := range perfume.Properties.Family {
		if _, err := tx.Exec(ctx, queries.InsertFamily, canonizedPerfume.Brand, canonizedPerfume.Name, perfume.Sex, family); err != nil {
			return err
		}
	}
	return nil
}

func updateNotes(ctx context.Context, tx pgx.Tx, query string, perfume models.Perfume, canonizedPerfume models.CanonizedPerfume, notes []string) error {
	for _, note := range notes {
		if _, err := tx.Exec(ctx, query, canonizedPerfume.Brand, canonizedPerfume.Name, perfume.Sex, note); err != nil {
			return err
		}
	}
	return nil
}

func updatePerfumeType(ctx context.Context, tx pgx.Tx, perfume models.Perfume, canonizedPerfume models.CanonizedPerfume) error {
	imageUrl := getPreferredImageUrl(perfume)
	if _, err := tx.Exec(ctx, queries.InsertPerfumeBaseInfo, canonizedPerfume.Brand, canonizedPerfume.Name, perfume.Sex, perfume.Brand, perfume.Name, perfume.Properties.Type, imageUrl); err != nil {
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
