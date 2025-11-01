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
		token := strings.TrimPrefix(r.Header.Get("Authorization"), prefix)
		if token != perfumeToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
