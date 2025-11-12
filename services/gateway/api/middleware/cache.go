package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/models/cache"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/models/perfume"
)

var (
	redisHost     = os.Getenv("REDIS_HOST")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	ttlEnv        = os.Getenv("TTL_SECONDS")
)

const defaultTTL = 1 * time.Hour

func Cache(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := getCacheKey(*r)
		ttl := getTTL()

		cacher := cache.GetOrCreateRedisCacher(redisHost, redisPort, redisPassword, ttl)
		cached, err := cacher.Load(r.Context(), key)
		if err == nil && cached != nil {
			suggestions, ok := cached.(perfume.Suggestions)
			if ok {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(suggestions); err != nil {
					log.Printf("Cannot encode cached response: %v\n", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}

		next(w, r)

		suggestionsValue, ok := r.Context().Value(cache.SuggestionsKey).(perfume.Suggestions)
		if ok {
			if err := cacher.Save(r.Context(), key, suggestionsValue); err != nil {
				log.Printf("Cannot cache: %v\n", err)
			}
		}
	}
}

func getCacheKey(r http.Request) string {
	brand := r.URL.Query().Get("brand")
	name := r.URL.Query().Get("name")
	sex := r.URL.Query().Get("sex")
	useAI := r.URL.Query().Get("use_ai")
	return fmt.Sprintf("%s:%s:%s:%s", brand, name, sex, useAI)
}

func getTTL() time.Duration {
	ttl, err := time.ParseDuration(ttlEnv)
	if err != nil {
		return defaultTTL
	}
	return ttl
}
