package middleware

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/errors"
)

func Cors(next http.HandlerFunc) http.HandlerFunc {
	allowedOrigins := []string{"http://frontend:3000", "http://localhost:3000"}

	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "" {
			if referer := r.Header.Get("Referer"); referer != "" {
				if u, err := url.Parse(referer); err == nil {
					origin = u.Scheme + "://" + u.Host
				}
			}
		}

		if !slices.Contains(allowedOrigins, origin) {
			gatewayErr := errors.ErrCORSNotAllowed(origin)
			gatewayErr.WriteHTTP(w)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}
