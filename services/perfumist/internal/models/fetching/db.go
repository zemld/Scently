package fetching

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

type perfumesFetchAndGlueResult struct {
	Perfumes []perfume.Perfume
	Status   int
}

type DB struct {
	url     string
	token   string
	timeout time.Duration
	client  *http.Client
}

func NewDB(url string, token string) *DB {
	return &DB{
		url:     url,
		token:   token,
		timeout: config.DBFetcherTimeout,
		client:  config.HTTPClient,
	}
}

func (f DB) Fetch(ctx context.Context, params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
	perfumesChan := make(chan perfumesFetchAndGlueResult, len(params))

	wg := sync.WaitGroup{}
	wg.Add(len(params))

	ctx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	for _, param := range params {
		go f.getPerfumesAsync(ctx, param, perfumesChan, &wg)
	}

	go func() {
		wg.Wait()
		close(perfumesChan)
	}()

	all, AllStatus := f.fetchPerfumeResults(ctx, perfumesChan)
	if AllStatus == http.StatusForbidden || AllStatus >= http.StatusInternalServerError {
		return nil, false
	}
	return all, true
}

func (f DB) getPerfumesAsync(ctx context.Context, params parameters.RequestPerfume, results chan<- perfumesFetchAndGlueResult, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	perfumes, status := f.getPerfumes(ctx, params)
	if status == http.StatusOK || status == http.StatusNotFound {
		results <- perfumesFetchAndGlueResult{Perfumes: perfumes, Status: status}
		return
	}
	results <- perfumesFetchAndGlueResult{Perfumes: nil, Status: status}
}

func (f DB) getPerfumes(ctx context.Context, p parameters.RequestPerfume) ([]perfume.Perfume, int) {
	r, err := http.NewRequestWithContext(ctx, "GET", f.url, nil)
	if err != nil {
		log.Printf("Can't create request: %v", err)
		return nil, http.StatusInternalServerError
	}
	p.AddToQuery(r)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", f.token))

	perfumeResponse, err := f.client.Do(r)
	if err != nil {
		log.Printf("Can't get perfumes: %v", err)
		return nil, http.StatusInternalServerError
	}
	defer perfumeResponse.Body.Close()

	body, err := io.ReadAll(perfumeResponse.Body)
	if err != nil {
		log.Printf("Can't read response body: %v", err)
		return nil, http.StatusInternalServerError
	}

	if perfumeResponse.StatusCode == http.StatusForbidden {
		log.Printf("Forbidden response: %s", string(body))
		return nil, http.StatusForbidden
	}

	if perfumeResponse.StatusCode == http.StatusNotFound {
		var perfumes perfume.PerfumeResponse
		if err := json.Unmarshal(body, &perfumes); err != nil {
			log.Printf("Can't unmarshal response: %v", err)
			return nil, http.StatusInternalServerError
		}
		log.Printf("Got %d perfumes (status: 404)", len(perfumes.Perfumes))
		return perfumes.Perfumes, http.StatusNotFound
	}

	if perfumeResponse.StatusCode == http.StatusInternalServerError {
		var perfumes perfume.PerfumeResponse
		if err := json.Unmarshal(body, &perfumes); err != nil {
			log.Printf("Can't unmarshal response: %v", err)
			return nil, http.StatusInternalServerError
		}
		log.Printf("Got %d perfumes (status: 500 - database error)", len(perfumes.Perfumes))
		return nil, http.StatusInternalServerError
	}

	if perfumeResponse.StatusCode == http.StatusOK {
		var perfumes perfume.PerfumeResponse
		if err := json.Unmarshal(body, &perfumes); err != nil {
			log.Printf("Can't unmarshal response: %v", err)
			return nil, http.StatusInternalServerError
		}
		log.Printf("Got %d perfumes", len(perfumes.Perfumes))
		return perfumes.Perfumes, http.StatusOK
	}

	return nil, http.StatusInternalServerError
}

func (f DB) fetchPerfumeResults(ctx context.Context, perfumesChan <-chan perfumesFetchAndGlueResult) ([]perfume.Perfume, int) {
	var perfumes []perfume.Perfume
	var status int

	for {
		select {
		case result, ok := <-perfumesChan:
			if !ok {
				return perfumes, status
			}
			if result.Status == http.StatusForbidden || result.Status == http.StatusInternalServerError {
				return perfumes, result.Status
			}
			if result.Status == http.StatusOK || result.Status == http.StatusNotFound {
				perfumes = append(perfumes, result.Perfumes...)
				if status == 0 || (status == http.StatusNotFound && result.Status == http.StatusOK) {
					status = result.Status
				}
			}
		case <-ctx.Done():
			return nil, http.StatusInternalServerError
		}
	}
}
