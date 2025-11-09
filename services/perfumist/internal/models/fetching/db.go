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

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

const (
	badServerStatus = http.StatusInternalServerError
)

type perfumesFetchAndGlueResult struct {
	Perfumes []perfume.Perfume
	Status   int
}

type Db struct {
	url     string
	token   string
	timeout time.Duration
}

func NewDb(url string, token string) *Db {
	return &Db{url: url, token: token, timeout: 2 * time.Second}
}

func (f Db) Fetch(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
	perfumesChan := make(chan perfumesFetchAndGlueResult, len(params))

	wg := sync.WaitGroup{}
	wg.Add(len(params))

	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	for _, param := range params {
		go f.getPerfumesAsync(ctx, param, perfumesChan, &wg)
	}

	go func() {
		wg.Wait()
		close(perfumesChan)
	}()

	all, AllStatus := f.fetchPerfumeResults(ctx, perfumesChan)
	if AllStatus >= badServerStatus || len(all) == 0 {
		return nil, false
	}
	return all, true
}

func (f Db) getPerfumesAsync(ctx context.Context, params parameters.RequestPerfume, results chan<- perfumesFetchAndGlueResult, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	perfumes, status := f.getPerfumes(ctx, params)
	if status != http.StatusOK {
		results <- perfumesFetchAndGlueResult{Perfumes: nil, Status: status}
		return
	}
	results <- perfumesFetchAndGlueResult{Perfumes: perfumes, Status: status}
}

func (f Db) getPerfumes(ctx context.Context, p parameters.RequestPerfume) ([]perfume.Perfume, int) {
	r, _ := http.NewRequestWithContext(ctx, "GET", f.url, nil)
	p.AddToQuery(r)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", f.token))

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

	var perfumes perfume.PerfumeResponse
	json.Unmarshal(body, &perfumes)
	log.Printf("Got %d perfumes", len(perfumes.Perfumes))
	if len(perfumes.Perfumes) == 0 {
		return perfumes.Perfumes, http.StatusNoContent
	}
	return perfumes.Perfumes, http.StatusOK
}

func (f Db) fetchPerfumeResults(ctx context.Context, perfumesChan <-chan perfumesFetchAndGlueResult) ([]perfume.Perfume, int) {
	var perfumes []perfume.Perfume
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
