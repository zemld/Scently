package application

import (
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/infrastructure"
)

func GetShortenResults(runResults []models.RunResults) []models.ShortenRunResults {
	runs := infrastructure.ReadResultsFromAllDirs("data/runs")
	shortenResults := make([]models.ShortenRunResults, 0, len(runs))
	for _, run := range runs {
		shortenResults = append(shortenResults, *models.NewShortenRunResults(run))
	}
	return shortenResults
}
