package infrastructure

import (
	"encoding/json"
	"log"
	"os"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/matching"
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
)

const threshold = 0.05

type testResultsWithWeights struct {
	matching.Weights `json:"weights"`
	Runs             []models.SingleSuggestResult `json:"runs"`
}

func SaveTestResults(outputPath string, inputs []models.Perfume, suggestions [][]models.Ranked, weights matching.Weights) {
	testResults := make([]models.SingleSuggestResult, 0, len(inputs))
	for i := range inputs {
		inputs[i].Properties.CalculateLeveledTags(threshold)
		inputs[i].Properties.CalculateLeveledCharacteristics()
		for j := range suggestions[i] {
			suggestions[i][j].Perfume.Properties.CalculateLeveledTags(threshold)
			suggestions[i][j].Perfume.Properties.CalculateLeveledCharacteristics()
		}
		testResults = append(testResults, models.SingleSuggestResult{
			Input:       inputs[i],
			Suggestions: suggestions[i],
		})
	}
	t := testResultsWithWeights{
		Weights: weights,
		Runs:    testResults,
	}
	encoded, err := json.Marshal(t)
	if err != nil {
		log.Fatalf("cannot marshal test results: %s", err)
	}
	os.WriteFile(outputPath, encoded, 0644)
}

func SaveShortenResults(outputPath string, shortenResults []models.ShortenRunResults) {
	encoded, err := json.Marshal(shortenResults)
	if err != nil {
		log.Fatalf("cannot marshal shorten results: %s", err)
	}
	os.WriteFile(outputPath, encoded, 0644)
}
