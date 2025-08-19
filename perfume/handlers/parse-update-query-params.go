package handlers

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/core"
)

func getUpdateParametersFromRequest(r *http.Request) *core.UpdateParameters {
	hard, err := strconv.ParseBool(r.URL.Query().Get("hard"))
	if err != nil {
		hard = false
	}
	p := core.NewUpdateParameters()
	if hard {
		setTruncateOptionIfNeeded(r, p)
	}
	return p
}

func setTruncateOptionIfNeeded(r *http.Request, p *core.UpdateParameters) {
	if getAndCheckPassword(r) {
		p.WithTruncate()
	}
}

func getAndCheckPassword(r *http.Request) bool {
	password := r.URL.Query().Get("password")
	return isPasswordValid(password)
}

func isPasswordValid(password string) bool {
	return password == getOriginalPassword()
}

func getOriginalPassword() string {
	password, _ := os.ReadFile(os.Getenv("HARD_UPDATE_PASSWORD_FILE"))
	return string(strings.Trim(string(password), "\n "))
}
