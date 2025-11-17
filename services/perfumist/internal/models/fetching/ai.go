package fetching

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

type aISuggestion struct {
	Perfumes []perfume.Perfume `json:"perfumes"`
}

type AI struct {
	url    string
	client *http.Client
}

func NewAI(url string) *AI {
	return &AI{
		url:    url,
		client: config.HTTPClient,
	}
}

func (f *AI) Fetch(ctx context.Context, params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
	ctx, cancel := context.WithTimeout(ctx, config.AIFetcherTimeout)
	defer cancel()

	if len(params) == 0 {
		return nil, false
	}

	r, err := http.NewRequestWithContext(ctx, "GET", f.url, nil)
	if err != nil {
		return nil, false
	}
	params[0].AddToQuery(r)

	response, err := f.client.Do(r)
	if err != nil {
		return nil, false
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, false
	}

	body, err := io.ReadAll(response.Body)
	if err != nil || len(body) == 0 {
		return nil, false
	}
	var suggestions aISuggestion
	if err := json.Unmarshal(body, &suggestions); err != nil {
		return nil, false
	}
	return suggestions.Perfumes, true
}
