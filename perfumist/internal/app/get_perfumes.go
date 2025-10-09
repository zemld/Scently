package app

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/util"
)

const (
	getPerfumesUrl = "http://perfume:8089/v1/perfumes/get"
)

func GetPerfumes(ctx context.Context, p util.GetParameters) ([]models.Perfume, int) {
	r, _ := http.NewRequestWithContext(ctx, "GET", getPerfumesUrl, nil)
	updateQuery(r, p)

	perfumeResponse, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Printf("Can't get perfumes: %v", err)
		return nil, http.StatusInternalServerError
	}
	defer perfumeResponse.Body.Close()

	if perfumeResponse.StatusCode >= 500 {
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
