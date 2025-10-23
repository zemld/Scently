package app

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
)

const (
	getPerfumesUrl = "http://perfume:8089/v1/perfumes/get"
)

const (
	badServerStatus = 500
)

func FetchPerfumes(ctx context.Context, params []parameters.RequestPerfume) ([]models.GluedPerfume, []models.GluedPerfume, int) {
	favouritePerfumesChan := make(chan perfumesFetchAndGlueResult)
	allPerfumesChan := make(chan perfumesFetchAndGlueResult, len(params)-1)

	wg := sync.WaitGroup{}
	wg.Add(len(params))

	go getAndGluePerfumesAsync(ctx, params[0], favouritePerfumesChan, &wg)
	for i := 1; i < len(params); i++ {
		go getAndGluePerfumesAsync(ctx, params[i], allPerfumesChan, &wg)
	}

	go func() {
		wg.Wait()
		close(favouritePerfumesChan)
		close(allPerfumesChan)
	}()

	fav, favStatus := fetchPerfumeResults(ctx, favouritePerfumesChan)
	all, AllStatus := fetchPerfumeResults(ctx, allPerfumesChan)
	if favStatus >= badServerStatus || AllStatus >= badServerStatus || len(all) == 0 {
		return nil, nil, http.StatusInternalServerError
	}
	if len(fav) == 0 {
		return nil, nil, http.StatusNoContent
	}
	return fav, all, http.StatusOK
}

type perfumesFetchAndGlueResult struct {
	Perfumes []models.GluedPerfume
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
	results <- perfumesFetchAndGlueResult{Perfumes: Glue(perfumes), Status: status}
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
}

func addQueryParameter(r *http.Request, key string, value string) {
	if value == "" {
		return
	}
	updatedQuery := r.URL.Query()
	updatedQuery.Add(key, value)
	r.URL.RawQuery = updatedQuery.Encode()
}

func fetchPerfumeResults(ctx context.Context, perfumesChan <-chan perfumesFetchAndGlueResult) ([]models.GluedPerfume, int) {
	var perfumes []models.GluedPerfume
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
