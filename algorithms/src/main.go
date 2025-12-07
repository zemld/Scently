package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/application"
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/matching"
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/infrastructure"
)

func main() {
	mode, alg, weights := application.ParseCLI(os.Args[1:])
	switch mode {
	case models.RunTests:
		all := application.ReadAndEnrichPerfumes()
		for i, weight := range weights {
			matcher := matching.GetMatcherByAlg(alg, weight)
			favs, results := application.RunTests(matcher, all)
			infrastructure.SaveTestResults(fmt.Sprintf("data/runs/%s/%s_%d.json", alg, alg, i), favs, results, weight)
		}
	case models.GetShortenResults:
		shortenResults := application.GetShortenResults(infrastructure.ReadResultsFromAllDirs("data/runs"))
		infrastructure.SaveShortenResults("data/shorten-results.json", shortenResults)
	default:
		log.Fatalf("unknown mode: %s", mode)
	}
}
