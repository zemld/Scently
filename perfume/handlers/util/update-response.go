package util

import (
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/models"
)

type UpdateResponse struct {
	SuccessfulPerfumes []models.Perfume    `json:"successful_perfumes"`
	FailedPerfumes     []models.Perfume    `json:"failed_perfumes"`
	State              core.ProcessedState `json:"state"`
}
