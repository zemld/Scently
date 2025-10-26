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

type UpdateStatus struct {
	SuccessfulPerfumes []models.Perfume `json:"successful_perfumes"`
	FailedPerfumes     []models.Perfume `json:"failed_perfumes"`
	State              ProcessedState   `json:"state"`
}

func NewUpdateStatus(success bool) *UpdateStatus {
	status := UpdateStatus{
		SuccessfulPerfumes: []models.Perfume{},
		FailedPerfumes:     []models.Perfume{},
		State:              NewProcessedState(),
	}
	status.State.Success = success
	return &status
}

func Update(params *UpdateParameters, perfumes []models.Perfume) UpdateStatus {
	config := config.NewConfig()
	ctx, cancel := CreateContext(config)
	defer cancel()

	conn, err := pgx.Connect(ctx, config.GetConnectionString())
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return *NewUpdateStatus(false)
	}
	defer conn.Close(ctx)

	tx, _ := conn.Begin(ctx)
	defer tx.Rollback(ctx)

	if params.IsTruncate {
		if !truncate(ctx, tx) {
			return *NewUpdateStatus(false)
		}
	}

	updateStatus := upsert(ctx, tx, perfumes)

	tx.Commit(ctx)
	return *updateStatus
}

func truncate(ctx context.Context, tx pgx.Tx) bool {
	_, err := tx.Exec(ctx, constants.Truncate)
	if err != nil {
		log.Printf("Error truncating tables: %v\n", err)
		return false
	}
	log.Println("Perfume tables truncated successfully")
	return true
}

func upsert(ctx context.Context, tx pgx.Tx, perfumes []models.Perfume) *UpdateStatus {
	updateStatus := NewUpdateStatus(true)
	for i, perfume := range perfumes {
		updateSavepointStatus(ctx, tx, constants.Savepoint, i)
		_, err := tx.Exec(ctx, constants.UpdatePerfumes, perfume.UnpackProperties()...)
		_, linkErr := tx.Exec(ctx, constants.UpdatePerfumeLinks, perfume.UnpackLinkedFields()...)
		if err != nil || linkErr != nil {
			log.Printf("Error updating perfume %s %s: %v\n, %v\n", perfume.Brand, perfume.Name, err, linkErr)
			updateSavepointStatus(ctx, tx, constants.RollbackSavepoint, i)
			updateStatus.FailedPerfumes = append(updateStatus.FailedPerfumes, perfume)
			updateStatus.State.FailedCount++
			continue
		}
		updateSavepointStatus(ctx, tx, constants.ReleaseSavepoint, i)
		updateStatus.SuccessfulPerfumes = append(updateStatus.SuccessfulPerfumes, perfume)
		updateStatus.State.SuccessfulCount++
	}

	return updateStatus
}

func updateSavepointStatus(ctx context.Context, tx pgx.Tx, cmd string, i int) {
	_, _ = tx.Exec(ctx, getSavepointQuery(cmd, i))
}

func getSavepointQuery(cmd string, i int) string {
	return fmt.Sprintf("%s%d", cmd, i)
}
