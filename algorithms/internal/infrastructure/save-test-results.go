package infrastructure

import (
	"encoding/json"
	"log"
	"os"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/matching"
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
)

type testResultsWithWeights struct {
	matching.Weights `json:"weights"`
	Runs             []testResult `json:"runs"`
}

type testResult struct {
	Input       models.Perfume  `json:"input"`
	Suggestions []models.Ranked `json:"suggestions"`
}

func SaveTestResults(outputPath string, inputs []models.Perfume, suggestions [][]models.Ranked, weights matching.Weights) {
	testResults := make([]testResult, 0, len(inputs))
	for i := range inputs {
		testResults = append(testResults, testResult{
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
