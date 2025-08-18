package responses

import "github.com/zemld/PerfumeRecommendationSystem/perfume/models"

type PerfumeResponse struct {
	Perfumes []models.Perfume `json:"perfumes"`
}
