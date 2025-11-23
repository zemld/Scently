package infrastructure

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"reflect"
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

func ReadNotesInfo[Number int | float64](path string) map[string]map[string]Number {
	header := make(map[int]string)

	csv := CsvIterator(path)

	headerLine, ok := <-csv
	if !ok {
		log.Fatal("cannot read header")
	}

	for i := 1; i < len(headerLine); i++ {
		header[i] = headerLine[i]
	}

	notes := make(map[string]map[string]Number)

	for noteLine := range csv {
		notes[noteLine[0]] = make(map[string]Number)
		for i := 1; i < len(noteLine); i++ {
			notes[noteLine[0]][header[i]] = parseNumber[Number](noteLine[i])
		}
	}

	return notes
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

func parseNumber[Number int | float64](value string) Number {
	var zeroValue Number
	isInt := reflect.TypeOf(zeroValue).Kind() == reflect.Int
	if isInt {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("cannot parse int: %s", err)
		}
		return Number(intVal)
	}
	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("cannot parse float: %s", err)
	}
	return Number(floatVal)
}
