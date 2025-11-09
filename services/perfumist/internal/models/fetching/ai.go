package fetching

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
)

const (
	aiSuggestUrl = "http://ai_advisor:8000/v1/advise"
)

type aISuggestion struct {
	Perfumes []perfume.Perfume `json:"perfumes"`
}

type AI struct {
	url string
}

func NewAI(url string) *AI {
	return &AI{url: url}
}

func (f *AI) Fetch(params []parameters.RequestPerfume) ([]perfume.Perfume, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, _ := http.NewRequestWithContext(ctx, "GET", f.url, nil)

	if len(params) == 0 {
		return nil, false
	}
	params[0].AddToQuery(r)

	response, err := http.DefaultClient.Do(r)
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
	err = json.Unmarshal(body, &suggestions)
	if err != nil {
		return nil, false
	}
	return suggestions.Perfumes, true
}
