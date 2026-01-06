package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/config"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/errors"
	"github.com/zemld/config-manager/pkg/cm"
)

func SuggestByTags(w http.ResponseWriter, r *http.Request) {
	if gatewayErr := validateTagsParameters(*r); gatewayErr != nil {
		gatewayErr.WriteHTTP(w)
		return
	}
	m := config.Manager()
	ctx, cancel := context.WithTimeout(r.Context(), getSuggestByTagsTimeout(m))
	defer cancel()

	perfumistUrl, err := getSuggestByTagsUrl(m)
	if err != nil {
		gatewayErr := errors.NewInternalError(err)
		gatewayErr.WriteHTTP(w)
		return
	}

	timeout := getSuggestByTagsTimeout(m)
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

func getSuggestByTagsUrl(cm cm.ConfigManager) (string, error) {
	return cm.GetString("suggest_by_tags_url")
}

func getSuggestByTagsTimeout(cm cm.ConfigManager) time.Duration {
	return cm.GetDurationWithDefault("non_ai_suggest_timeout", 8*time.Second)
}

func validateTagsParameters(r http.Request) *errors.GatewayError {
	tags := r.URL.Query().Get("tags")
	if tags == "" {
		return errors.ErrBadRequest(fmt.Errorf("tags is required"))
	}
	return nil
}
