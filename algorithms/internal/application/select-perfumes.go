package application

import "github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"

func SelectConcretePerfume(brand string, name string, sex string, perfumes []models.Perfume) []models.Perfume {
	selectedPerfumes := make([]models.Perfume, 0, len(perfumes))
	for _, perfume := range perfumes {
		if perfume.Brand == brand && perfume.Name == name && perfume.Sex == sex {
			selectedPerfumes = append(selectedPerfumes, perfume)
		}
	}
	return selectedPerfumes
}

func SelectPerfumesBySex(sex string, perfumes []models.Perfume) []models.Perfume {
	selectedPerfumes := make([]models.Perfume, 0, len(perfumes))
	for _, perfume := range perfumes {
		if perfume.Sex == sex || perfume.Sex == "unisex" {
			selectedPerfumes = append(selectedPerfumes, perfume)
		}
	}
	return selectedPerfumes
}
