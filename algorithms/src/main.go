package main

import (
	"fmt"
	"os"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/application"
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/matching"
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/infrastructure"
)

func main() {
	all := application.ReadAndEnrichPerfumes()
	weights := application.ParseCLIAndGetWeights(os.Args[1:])
	for i, weight := range weights {
		matcher := matching.GetMatcherByAlg(matching.AlgType(os.Args[1]), weight)
		favs, results := application.RunTests(matcher, all)
		infrastructure.SaveTestResults(fmt.Sprintf("data/runs/%s/%s_%d.json", os.Args[1], os.Args[1], i), favs, results, weight)
	}

}
