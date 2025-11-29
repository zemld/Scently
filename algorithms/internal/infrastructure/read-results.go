package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
)

func ReadResultsFromAllDirs(path string) []models.RunResults {
	dirs, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("cannot read directory %s: %s", path, err)
	}
	runResults := make([]models.RunResults, 0, len(dirs))
	for dirID, dir := range dirs {
		dirRunResults := ReadDirResults(fmt.Sprintf("%s/%s", path, dir.Name()))
		for i := range dirRunResults {
			currentRunID := dirRunResults[i].ID
			updatedRunID := fmt.Sprintf("%d-%s", dirID, currentRunID)
			dirRunResults[i].ID = updatedRunID
		}
		runResults = append(runResults, dirRunResults...)
	}
	return runResults
}

func ReadDirResults(path string) []models.RunResults {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("cannot read directory %s: %s", path, err)
	}
	runResults := make([]models.RunResults, 0, len(files))
	for fileID, file := range files {
		singleRunResults := ReadFileResults(fmt.Sprintf("%s/%s", path, file.Name()))
		runResults = append(runResults, models.RunResults{
			ID:   fmt.Sprintf("%d", fileID),
			Runs: singleRunResults,
		})
	}
	return runResults
}

func ReadFileResults(path string) []models.SingleSuggestResult {
	contents, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("cannot read file %s: %s", path, err)
	}
	var testResults testResultsWithWeights
	if err := json.Unmarshal(contents, &testResults); err != nil {
		log.Fatalf("cannot unmarshal test results: %s", err)
	}
	runResults := make([]models.SingleSuggestResult, 0, len(testResults.Runs))
	for _, run := range testResults.Runs {
		runResults = append(runResults, models.SingleSuggestResult{
			Input:       run.Input,
			Suggestions: run.Suggestions,
		})
	}
	return runResults
}
