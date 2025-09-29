package internal

import (
	"fmt"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/models"
)

func Glue(perfumes []models.Perfume) []models.GluedPerfume {
	return fetchGluedPerfumesFromMap(glue(perfumes))
}

func glue(perfumes []models.Perfume) map[string]models.GluedPerfume {
	gluedMap := make(map[string]models.GluedPerfume)
	for _, perfume := range perfumes {
		key := getKey(perfume)
		if _, ok := gluedMap[key]; ok {
			gluedMap[key].Links[perfume.Volume] = perfume.Link
		} else {
			gluedMap[key] = models.NewGluedPerfume(perfume)
		}
	}
	return gluedMap
}

func getKey(perfume models.Perfume) string {
	return fmt.Sprintf("%s%s", perfume.Brand, perfume.Name)
}

func fetchGluedPerfumesFromMap(gluedMap map[string]models.GluedPerfume) []models.GluedPerfume {
	glued := make([]models.GluedPerfume, len(gluedMap))
	var i uint64 = 0
	for _, perfume := range gluedMap {
		glued[i] = perfume
		i++
	}
	return glued
}
