package middleware

import (
	"net/http"
	"os"
	"strings"
)

var perfumeToken = os.Getenv("PERFUME_INTERNAL_TOKEN")

const prefix = "Bearer "

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawToken := r.Header.Get("Authorization")
		if !strings.HasPrefix(rawToken, prefix) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(rawToken, prefix)
		if token != perfumeToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
