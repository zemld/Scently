package util

import (
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/models"
)

type glueKey struct {
	Brand string
	Name  string
}

func GluePerfumes(perfumes []models.Perfume) []models.GluedPerfume {
	glued := make(map[glueKey]models.GluedPerfume)

	for _, perfume := range perfumes {
		key := glueKey{Brand: perfume.Brand, Name: perfume.Name}
		p, found := glued[key]
		if found {
			p.Links[perfume.Volume] = perfume.Link
		} else {
			glued[key] = models.NewGluedPerfume(perfume)
		}
	}
	result := make([]models.GluedPerfume, 0, len(glued))
	for _, p := range glued {
		result = append(result, p)
	}
	return result
}
