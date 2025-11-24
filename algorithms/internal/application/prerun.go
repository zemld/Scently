package application

import (
	"log"
	"maps"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/infrastructure"
)

const (
	perfumesPath        = "data/all_perfumes_merged.json"
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

	for i := range perfumes {
		perfumes[i].Properties.EnrichedUpperNotes = enrichNotes(perfumes[i].Properties.UpperNotes, tags, characteristics)
		perfumes[i].Properties.EnrichedCoreNotes = enrichNotes(perfumes[i].Properties.CoreNotes, tags, characteristics)
		perfumes[i].Properties.EnrichedBaseNotes = enrichNotes(perfumes[i].Properties.BaseNotes, tags, characteristics)
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

		if noteTags, ok := tags[note]; ok && noteTags != nil {
			maps.Copy(enrichedNote.Tags, noteTags)
		}
		if noteCharacteristics, ok := characteristics[note]; ok && noteCharacteristics != nil {
			maps.Copy(enrichedNote.Characteristics, noteCharacteristics)
		}

		enrichedNotes = append(enrichedNotes, enrichedNote)
	}

	return enrichedNotes
}
