package infrastructure

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
)

const (
	readFileMsg  = "cannot read file %s with error: %s"
	unmarshalMsg = "cannot unmarshal data: %s"
	readCsvMsg   = "cannot read record: %s"
)

func ReadPerfumes(path string) []models.Perfume {
	contents, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf(readFileMsg, path, err)
	}
	var perfumes []models.Perfume
	if err := json.Unmarshal(contents, &perfumes); err != nil {
		log.Fatalf(unmarshalMsg, err)
	}
	return perfumes
}

func ReadTags(path string, notes map[string]models.EnrichedNote) {
	header := make(map[int]string)

	csv := CsvIterator(path)

	headerLine, ok := <-csv
	if !ok {
		log.Fatal("cannot read header")
	}

	for i := 1; i < len(headerLine); i++ {
		header[i] = headerLine[i]
	}

	for noteLine := range csv {
		notes[noteLine[0]] = models.EnrichedNote{
			Name: noteLine[0],
			Tags: make(map[string]int),
		}

		for i := 1; i < len(noteLine); i++ {
			count, err := strconv.Atoi(noteLine[i])
			if err != nil {
				log.Panicf("cannot convert %s to int: %s", noteLine[i], err)
				notes[noteLine[0]].Tags[header[i]] = 0
			}
			notes[noteLine[0]].Tags[header[i]] = count
		}
	}
}

func CsvIterator(path string) <-chan []string {
	ch := make(chan []string)

	go func() {
		defer close(ch)

		contents, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf(readFileMsg, path, err)
			return
		}

		reader := csv.NewReader(bytes.NewBuffer(contents))
		reader.Comma = ';'

		for {
			record, err := reader.Read()
			if err != nil {
				return
			}
			ch <- record
		}
	}()

	return ch
}
