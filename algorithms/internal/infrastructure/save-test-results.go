package infrastructure

import (
	"encoding/json"
	"log"
	"os"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
)

type TestResult struct {
	Input       models.Perfume  `json:"input"`
	Suggestions []models.Ranked `json:"suggestions"`
}

func SaveTestResults(outputPath string, inputs []models.Perfume, suggestions [][]models.Ranked) {
	testResults := make([]TestResult, 0, len(inputs))
	for i := range inputs {
		testResults = append(testResults, TestResult{
			Input:       inputs[i],
			Suggestions: suggestions[i],
		})
	}
	encoded, err := json.Marshal(testResults)
	if err != nil {
		log.Fatalf("cannot marshal test results: %s", err)
	}
	os.WriteFile(outputPath, encoded, 0644)
}
