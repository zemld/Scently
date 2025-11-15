package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/models/cache"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/models/canonization"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/models/perfume"
)

var (
	redisHost     = os.Getenv("REDIS_HOST")
	redisPort     = os.Getenv("REDIS_PORT")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	ttlEnv        = os.Getenv("TTL_SECONDS")
)

const defaultTTL = 1 * time.Hour

type responseWriter struct {
	http.ResponseWriter
	body       []byte
	statusCode int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = append(rw.body, b...)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func Cache(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := getCacheKey(*r)
		ttl := getTTL()

		cacher := cache.GetOrCreateRedisCacher(redisHost, redisPort, redisPassword, ttl)
		cached, err := cacher.Load(r.Context(), key)
		if err == nil && cached != nil {
			var suggestions perfume.Suggestions
			if err := json.Unmarshal(cached, &suggestions); err == nil {
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(suggestions); err != nil {
					log.Printf("Cannot encode cached response: %v\n", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				log.Printf("Cache hit for key: %s\n", key)
				log.Printf("Cached response body: %s\n", string(cached))
				return
			}
		}

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next(rw, r)

		if rw.statusCode == http.StatusOK && len(rw.body) > 0 {
			if err := cacher.Save(r.Context(), key, rw.body); err != nil {
				log.Printf("Cannot cache: %v\n", err)
			} else {
				log.Printf("Cached response for key: %s\n", key)
				log.Printf("Cached response body: %s\n", string(rw.body))
			}
		}
	}
}

func getCacheKey(r http.Request) string {
	canonizer := canonization.DefaultCanonizer{}
	keys := []string{
		r.URL.Query().Get("brand"),
		r.URL.Query().Get("name"),
		r.URL.Query().Get("sex"),
		r.URL.Query().Get("use_ai"),
	}
	return canonizer.Canonize(keys)
}

func getTTL() time.Duration {
	ttl, err := time.ParseDuration(os.Getenv(ttlEnv))
	if err != nil {
		return defaultTTL
	}
	return ttl
}
