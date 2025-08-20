package util

import (
	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/core"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/models"
)

type PerfumeResponse struct {
	Perfumes []models.Perfume    `json:"perfumes"`
	State    core.ProcessedState `json:"state"`
}
type PerfumeCollection struct {
	Perfumes []models.Perfume `json:"perfumes"`
}
