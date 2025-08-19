package responses

import "github.com/zemld/PerfumeRecommendationSystem/perfume/models"

type PerfumeCollection struct {
	Perfumes []models.Perfume `json:"perfumes"`
}
