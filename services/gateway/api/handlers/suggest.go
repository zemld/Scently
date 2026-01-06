package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/config"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/errors"
	"github.com/zemld/config-manager/pkg/cm"
)

func Suggest(w http.ResponseWriter, r *http.Request) {
	if gatewayErr := validateParameters(*r); gatewayErr != nil {
		gatewayErr.WriteHTTP(w)
		return
	}
	m := config.Manager()
	ctx, cancel := context.WithTimeout(r.Context(), getTimeoutFromRequest(*r, m))
	defer cancel()

	perfumistUrl, err := getSuggestionUrlFromRequest(*r, m)
	if err != nil {
		gatewayErr := errors.NewInternalError(err)
		gatewayErr.WriteHTTP(w)
		return
	}

	timeout := getTimeoutFromRequest(*r, m)
	resp, body, err := proxyRequestToPerfumist(ctx, perfumistUrl, r, timeout, true)
	if err != nil {
		gatewayErr := errors.NewInternalError(err)
		gatewayErr.WriteHTTP(w)
		return
	}
	defer resp.Body.Close()

	if err := handlePerfumistResponse(w, resp, body); err != nil {
		log.Printf("Error handling perfumist response: %v\n", err)
	}
}

func getSuggestionUrlFromRequest(r http.Request, cm cm.ConfigManager) (string, error) {
	if strings.EqualFold(r.URL.Query().Get("use_ai"), "true") {
		return cm.GetString("ai_suggest_url")
	}
	return cm.GetString("suggest_url")
}

func getTimeoutFromRequest(r http.Request, cm cm.ConfigManager) time.Duration {
	if strings.EqualFold(r.URL.Query().Get("use_ai"), "true") {
		return cm.GetDurationWithDefault("ai_suggest_timeout", 20*time.Second)
	}
	return cm.GetDurationWithDefault("non_ai_suggest_timeout", 8*time.Second)
}

func validateParameters(r http.Request) *errors.GatewayError {
	if r.URL.Query().Get("brand") == "" {
		return errors.ErrBadRequest(fmt.Errorf("brand is required"))
	}
	if r.URL.Query().Get("name") == "" {
		return errors.ErrBadRequest(fmt.Errorf("name is required"))
	}
	return nil
}
