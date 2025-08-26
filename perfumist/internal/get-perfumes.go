package internal

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/models"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/util"
)

const (
	getPerfumesUrl = "http://perfume:8089/v1/perfumes/get"
)

func GetPerfumes(p util.GetParameters) ([]models.Perfume, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, _ := http.NewRequestWithContext(ctx, "GET", getPerfumesUrl, nil)
	updateQuery(r, p)

	perfumeResponse, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Printf("Can't get perfumes: %v", err)
		return nil, false
	}
	defer perfumeResponse.Body.Close()

	if perfumeResponse.StatusCode >= 500 {
		log.Printf("Bad response status: %v", perfumeResponse.Status)
		return nil, false
	}
	body, err := io.ReadAll(perfumeResponse.Body)
	if err != nil {
		log.Printf("Can't read response body: %v", err)
		return nil, false
	}

	var perfumes models.PerfumeResponse
	json.Unmarshal(body, &perfumes)
	log.Printf("Got %d perfumes", len(perfumes.Perfumes))
	return perfumes.Perfumes, true
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
