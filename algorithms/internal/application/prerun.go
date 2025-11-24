package application

import (
	"log"
	"maps"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/infrastructure"
)

const (
	perfumesPath        = "data/all_perfumes.json"
	tagsPath            = "data/note_tags_dataset_filled.csv"
	characteristicsPath = "data/note_characteristics_filled.csv"
)

func ReadAndEnrichPerfumes() []models.Perfume {
	perfumes := infrastructure.ReadPerfumes(perfumesPath)
	if perfumes == nil {
		log.Fatal("cannot read perfumes")
	}
	tags := infrastructure.ReadNotesInfo[int](tagsPath)
	if tags == nil {
		log.Fatal("cannot read tags")
	}
	characteristics := infrastructure.ReadNotesInfo[float64](characteristicsPath)
	if characteristics == nil {
		log.Fatal("cannot read characteristics")
	}

	for _, perfume := range perfumes {
		perfume.Properties.EnrichedUpperNotes = enrichNotes(perfume.Properties.UpperNotes, tags, characteristics)
		perfume.Properties.EnrichedCoreNotes = enrichNotes(perfume.Properties.CoreNotes, tags, characteristics)
		perfume.Properties.EnrichedBaseNotes = enrichNotes(perfume.Properties.BaseNotes, tags, characteristics)
	}
	return perfumes
}

func enrichNotes(notes []string, tags map[string]map[string]int, characteristics map[string]map[string]float64) []models.EnrichedNote {
	enrichedNotes := make([]models.EnrichedNote, 0, len(notes))
	for _, note := range notes {
		enrichedNote := models.EnrichedNote{
			Name:            note,
			Tags:            make(map[string]int),
			Characteristics: make(map[string]float64),
		}

		maps.Copy(enrichedNote.Tags, tags[note])
		maps.Copy(enrichedNote.Characteristics, characteristics[note])

		enrichedNotes = append(enrichedNotes, enrichedNote)
	}

	return enrichedNotes
}
