package middleware

import (
	"context"
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

		cacher, err := cache.NewRedisCacher(redisHost, redisPort, redisPassword, ttl)
		if err != nil {
			log.Printf("Cannot create Redis cacher: %v\n", err)
		}

		if tryLoadFromCache(r.Context(), cacher, key, w) {
			return
		}

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next(rw, r)

		if rw.statusCode == http.StatusOK && len(rw.body) > 0 && cacher != nil {
			if err := cacher.Save(r.Context(), key, rw.body); err != nil {
				log.Printf("Cannot cache: %v\n", err)
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

func tryLoadFromCache(ctx context.Context, cacher cache.Loader, key string, w http.ResponseWriter) bool {
	cached, err := cacher.Load(ctx, key)
	if err != nil || cached == nil {
		return false
	}
	var suggestions perfume.Suggestions
	if err := json.Unmarshal(cached, &suggestions); err != nil {
		return false
	}
	if len(suggestions.Perfumes) == 0 {
		return false
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(suggestions); err != nil {
		return false
	}
	return true
}
