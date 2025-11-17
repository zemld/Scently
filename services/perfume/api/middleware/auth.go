package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/errors"
)

var perfumeToken = os.Getenv("PERFUME_INTERNAL_TOKEN")

const prefix = "Bearer "

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawToken := r.Header.Get("Authorization")
		if !strings.HasPrefix(rawToken, prefix) {
			authErr := errors.NewAuthError("missing or invalid authorization header")
			handleAuthError(w, authErr)
			return
		}
		token := strings.TrimPrefix(rawToken, prefix)
		if token != perfumeToken {
			authErr := errors.NewAuthError("invalid token")
			handleAuthError(w, authErr)
			return
		}
		next(w, r)
	}
}

func handleAuthError(w http.ResponseWriter, err *errors.AuthError) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(err.HTTPStatus())
	w.Write([]byte(err.Error()))
}
