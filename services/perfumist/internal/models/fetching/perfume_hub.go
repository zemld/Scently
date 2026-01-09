package fetching

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zemld/Scently/models"
	"github.com/zemld/Scently/perfumist/internal/models/parameters"
	"github.com/zemld/Scently/perfumist/internal/models/perfume"
	"github.com/zemld/config-manager/pkg/cm"
)

type perfumesFetchAndGlueResult struct {
	Perfumes []models.Perfume
	Status   int
}

type PerfumeHub struct {
	url     string
	token   string
	timeout time.Duration
	client  *http.Client
	cm      cm.ConfigManager
}

func NewPerfumeHub(url string, token string, cm cm.ConfigManager) *PerfumeHub {
	return &PerfumeHub{
		url:     url,
		token:   token,
		timeout: cm.GetDurationWithDefault("perfume_hub_fetcher_timeout", 5*time.Second),
		client:  http.DefaultClient,
		cm:      cm,
	}
}

func (f *PerfumeHub) FetchMany(ctx context.Context, params []parameters.RequestPerfume) <-chan models.Perfume {
	allPerfumesChan := make(chan models.Perfume)

	wg := sync.WaitGroup{}
	wg.Add(len(params))
	for _, param := range params {
		go func(p parameters.RequestPerfume) {
			defer wg.Done()
			perfumesChan := f.Fetch(ctx, p)
			for {
				select {
				case <-ctx.Done():
					return
				case perfume, ok := <-perfumesChan:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case allPerfumesChan <- perfume:
					}
				}
			}
		}(param)
	}

	go func() {
		wg.Wait()
		close(allPerfumesChan)
	}()
	return allPerfumesChan
}

func (f *PerfumeHub) Fetch(ctx context.Context, parameter parameters.RequestPerfume) <-chan models.Perfume {
	if parameter.Brand != "" || parameter.Name != "" {
		return f.fetchConcretePerfume(ctx, parameter)
	}
	return f.fetchAllPerfumes(ctx, parameter)
}

func (f *PerfumeHub) fetchConcretePerfume(ctx context.Context, parameter parameters.RequestPerfume) <-chan models.Perfume {
	perfumeChan := make(chan models.Perfume)
	go func() {
		defer close(perfumeChan)
		perfumes, status := f.getPerfumes(ctx, parameter)
		if status == http.StatusOK {
			for _, perfume := range perfumes {
				select {
				case <-ctx.Done():
					return
				case perfumeChan <- perfume:
				}
			}
		}
	}()
	return perfumeChan
}

func (f *PerfumeHub) fetchAllPerfumes(ctx context.Context, parameter parameters.RequestPerfume) <-chan models.Perfume {
	var pageNumber atomic.Uint32
	pageNumber.Store(1)

	perfumeChan := make(chan models.Perfume)
	workersCount := f.cm.GetIntWithDefault("threads_count", 8)
	wg := sync.WaitGroup{}
	wg.Add(workersCount)
	for range workersCount {
		go func() {
			defer wg.Done()
			localParameter := parameter
			for {
				localParameter.Page = pageNumber.Load()
				pageNumber.Add(1)
				perfumes, status := f.getPerfumes(ctx, localParameter)
				if status != http.StatusOK {
					break
				}
				for _, perfume := range perfumes {
					select {
					case <-ctx.Done():
						return
					case perfumeChan <- perfume:
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(perfumeChan)
	}()
	return perfumeChan
}

func (f *PerfumeHub) getPerfumes(ctx context.Context, p parameters.RequestPerfume) ([]models.Perfume, int) {
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
