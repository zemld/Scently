package internal

import (
	"context"
	"net/http"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/util"
)

const (
	getPerfumesUrl = "http://db-api:8089/v1/perfumes/get"
)

func GetPerfumes(p util.GetParameters) []models.Perfume {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, _ := http.NewRequestWithContext(ctx, "GET", getPerfumesUrl, nil)
	updateQuery(r, p)

	// TODO: add response parsing
	perfumeResponse, _ := http.DefaultClient.Do(r)
	defer perfumeResponse.Body.Close()

	return []models.Perfume{}
}

func updateQuery(r *http.Request, p util.GetParameters) {
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
