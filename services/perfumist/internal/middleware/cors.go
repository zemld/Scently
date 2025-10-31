package middleware

import (
	"net/http"
	"slices"
)

func Cors(next http.HandlerFunc) http.HandlerFunc {
	allowedOrigins := []string{"http://frontend:3000"}
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if !slices.Contains(allowedOrigins, origin) {
			http.Error(w, "CORS not allowed", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
