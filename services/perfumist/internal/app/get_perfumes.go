package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
)

const (
	getPerfumesUrl          = "http://perfume:8000/v1/perfumes/get"
	perfumeInternalTokenEnv = "PERFUME_INTERNAL_TOKEN"
)

const (
	badServerStatus = 500
)

func FetchPerfumes(ctx context.Context, params []parameters.RequestPerfume) ([]models.Perfume, bool) {
	perfumesChan := make(chan perfumesFetchAndGlueResult, len(params))

	wg := sync.WaitGroup{}
	wg.Add(len(params))

	for _, param := range params {
		go getAndGluePerfumesAsync(ctx, param, perfumesChan, &wg)
	}

	go func() {
		wg.Wait()
		close(perfumesChan)
	}()

	all, AllStatus := fetchPerfumeResults(ctx, perfumesChan)
	if AllStatus >= badServerStatus || len(all) == 0 {
		return nil, false
	}
	return all, true
}

type perfumesFetchAndGlueResult struct {
	Perfumes []models.Perfume
	Status   int
}

func getAndGluePerfumesAsync(ctx context.Context, params parameters.RequestPerfume, results chan<- perfumesFetchAndGlueResult, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	perfumes, status := getPerfumes(ctx, params)
	if status != http.StatusOK {
		results <- perfumesFetchAndGlueResult{Perfumes: nil, Status: status}
		return
	}
	results <- perfumesFetchAndGlueResult{Perfumes: perfumes, Status: status}
}

func getPerfumes(ctx context.Context, p parameters.RequestPerfume) ([]models.Perfume, int) {
	r, _ := http.NewRequestWithContext(ctx, "GET", getPerfumesUrl, nil)
	updateQuery(r, p)

	perfumeResponse, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Printf("Can't get perfumes: %v", err)
		return nil, http.StatusInternalServerError
	}
	defer perfumeResponse.Body.Close()

	if perfumeResponse.StatusCode >= badServerStatus {
		log.Printf("Bad response status: %v", perfumeResponse.Status)
		return nil, perfumeResponse.StatusCode
	}
	body, err := io.ReadAll(perfumeResponse.Body)
	if err != nil {
		log.Printf("Can't read response body: %v", err)
		return nil, http.StatusInternalServerError
	}

	var perfumes models.PerfumeResponse
	json.Unmarshal(body, &perfumes)
	log.Printf("Got %d perfumes", len(perfumes.Perfumes))
	if len(perfumes.Perfumes) == 0 {
		return perfumes.Perfumes, http.StatusNoContent
	}
	return perfumes.Perfumes, http.StatusOK
}

func updateQuery(r *http.Request, p parameters.RequestPerfume) {
	addQueryParameter(r, "brand", p.Brand)
	addQueryParameter(r, "name", p.Name)
	if p.Sex == "male" || p.Sex == "female" {
		addQueryParameter(r, "sex", p.Sex)
	}
	log.Printf("Perfume internal token: %s", os.Getenv(perfumeInternalTokenEnv))
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv(perfumeInternalTokenEnv)))
}

func addQueryParameter(r *http.Request, key string, value string) {
	if value == "" {
		return
	}
	updatedQuery := r.URL.Query()
	updatedQuery.Set(key, value)
	r.URL.RawQuery = updatedQuery.Encode()
}

func fetchPerfumeResults(ctx context.Context, perfumesChan <-chan perfumesFetchAndGlueResult) ([]models.Perfume, int) {
	var perfumes []models.Perfume
	var status int

	for {
		select {
		case result, ok := <-perfumesChan:
			if !ok {
				return perfumes, status
			}
			if result.Status >= badServerStatus {
				return perfumes, result.Status
			}
			if result.Status == http.StatusOK {
				perfumes = append(perfumes, result.Perfumes...)
			}
		case <-ctx.Done():
			return nil, http.StatusInternalServerError
		}
	}
}
