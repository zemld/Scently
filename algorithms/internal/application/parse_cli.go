package application

import (
	"log"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/matching"
)

func ParseCLIAndGetWeights(args []string) []matching.Weights {
	if len(args) == 0 {
		log.Fatal("no important arg of alg type")
	}
	alg := args[0]
	weights := getWeights(matching.AlgType(alg))

	if weights == nil {
		log.Fatalf("unknown algorithm: %s", alg)
		return nil
	}

	return weights
}

func getWeights(alg matching.AlgType) []matching.Weights {
	switch alg {
	case matching.OverlayAlg:
		return []matching.Weights{
			*matching.NewOverlayWeights(0.15, 0.45, 0.4, 0.4, 0.55, 0.05),
			*matching.NewOverlayWeights(0.15, 0.45, 0.4, 0.3, 0.65, 0.05),
			*matching.NewOverlayWeights(0.15, 0.45, 0.4, 0.2, 0.75, 0.05),

			*matching.NewOverlayWeights(0.25, 0.45, 0.3, 0.4, 0.55, 0.05),
			*matching.NewOverlayWeights(0.25, 0.45, 0.3, 0.3, 0.65, 0.05),
			*matching.NewOverlayWeights(0.25, 0.45, 0.3, 0.2, 0.75, 0.05),

			*matching.NewOverlayWeights(0.2, 0.35, 0.45, 0.4, 0.55, 0.05),
			*matching.NewOverlayWeights(0.2, 0.35, 0.45, 0.3, 0.65, 0.05),
			*matching.NewOverlayWeights(0.2, 0.35, 0.45, 0.2, 0.75, 0.05),
		}
	case matching.TagsOverlayAlg, matching.TagsAlg, matching.CharacteristicsAlg:
		return []matching.Weights{
			*matching.NewBaseWeights(0.15, 0.45, 0.4),
			*matching.NewBaseWeights(0.25, 0.45, 0.3),
			*matching.NewBaseWeights(0.2, 0.35, 0.45),
		}
	case matching.SmartAlg:
		return []matching.Weights{
			*matching.NewSmartWeights(0.15, 0.45, 0.4, 0.2, 0.8),
			*matching.NewSmartWeights(0.15, 0.45, 0.4, 0.3, 0.7),
			*matching.NewSmartWeights(0.15, 0.45, 0.4, 0.4, 0.6),
			*matching.NewSmartWeights(0.15, 0.45, 0.4, 0.5, 0.5),
			*matching.NewSmartWeights(0.15, 0.45, 0.4, 0.6, 0.4),
			*matching.NewSmartWeights(0.15, 0.45, 0.4, 0.7, 0.3),
			*matching.NewSmartWeights(0.15, 0.45, 0.4, 0.8, 0.2),

			*matching.NewSmartWeights(0.25, 0.45, 0.3, 0.2, 0.8),
			*matching.NewSmartWeights(0.25, 0.45, 0.3, 0.3, 0.7),
			*matching.NewSmartWeights(0.25, 0.45, 0.3, 0.4, 0.6),
			*matching.NewSmartWeights(0.25, 0.45, 0.3, 0.5, 0.5),
			*matching.NewSmartWeights(0.25, 0.45, 0.3, 0.6, 0.4),
			*matching.NewSmartWeights(0.25, 0.45, 0.3, 0.7, 0.3),
			*matching.NewSmartWeights(0.25, 0.45, 0.3, 0.8, 0.2),

			*matching.NewSmartWeights(0.2, 0.35, 0.45, 0.2, 0.8),
			*matching.NewSmartWeights(0.2, 0.35, 0.45, 0.3, 0.7),
			*matching.NewSmartWeights(0.2, 0.35, 0.45, 0.4, 0.6),
			*matching.NewSmartWeights(0.2, 0.35, 0.45, 0.5, 0.5),
			*matching.NewSmartWeights(0.2, 0.35, 0.45, 0.6, 0.4),
			*matching.NewSmartWeights(0.2, 0.35, 0.45, 0.7, 0.3),
			*matching.NewSmartWeights(0.2, 0.35, 0.45, 0.8, 0.2),
		}
	case matching.SmartEnhancedAlg:
		return []matching.Weights{
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.3, 0.6, 0.1),
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.4, 0.5, 0.1),
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.5, 0.4, 0.1),
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.6, 0.3, 0.1),
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.2, 0.6, 0.2),
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.3, 0.5, 0.2),
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.4, 0.4, 0.2),
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.5, 0.3, 0.2),
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.6, 0.2, 0.2),
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.2, 0.5, 0.3),
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.3, 0.4, 0.3),
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.4, 0.3, 0.3),
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.2, 0.4, 0.4),
			*matching.NewSmartEnhancedWeights(0.15, 0.45, 0.4, 0.3, 0.3, 0.4),

			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.3, 0.6, 0.1),
			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.4, 0.5, 0.1),
			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.5, 0.4, 0.1),
			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.6, 0.3, 0.1),
			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.2, 0.6, 0.2),
			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.3, 0.5, 0.2),
			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.4, 0.4, 0.2),
			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.5, 0.3, 0.2),
			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.6, 0.2, 0.2),
			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.2, 0.5, 0.3),
			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.3, 0.4, 0.3),
			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.4, 0.3, 0.3),
			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.2, 0.4, 0.4),
			*matching.NewSmartEnhancedWeights(0.25, 0.45, 0.3, 0.3, 0.3, 0.4),

			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.3, 0.6, 0.1),
			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.4, 0.5, 0.1),
			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.5, 0.4, 0.1),
			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.6, 0.3, 0.1),
			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.2, 0.6, 0.2),
			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.3, 0.5, 0.2),
			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.4, 0.4, 0.2),
			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.5, 0.3, 0.2),
			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.6, 0.2, 0.2),
			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.2, 0.5, 0.3),
			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.3, 0.4, 0.3),
			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.4, 0.3, 0.3),
			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.2, 0.4, 0.4),
			*matching.NewSmartEnhancedWeights(0.2, 0.35, 0.45, 0.3, 0.3, 0.4),
		}
	default:
		log.Fatalf("unknown algorithm: %s", alg)
		return nil
	}
}
