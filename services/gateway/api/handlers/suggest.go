package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/errors"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/models/perfume"
)

const (
	aITimeout    = "SUGGEST_AI_TIMEOUT"
	nonAITimeout = "SUGGEST_TIMEOUT"
)

const (
	defaultAITimeout    = 20 * time.Second
	defaultNonAITimeout = 2 * time.Second
)

var (
	suggestUrl = os.Getenv("PERFUMIST_URL")
)

func Suggest(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), getTimeoutFromRequest(*r))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, suggestUrl, nil)
	if err != nil {
		gatewayErr := errors.NewInternalError(err)
		gatewayErr.WriteHTTP(w)
		return
	}
	req.URL.RawQuery = r.URL.Query().Encode()

	timeout := getTimeoutFromRequest(*r)
	client := getHTTPClient(timeout)
	resp, err := client.Do(req)
	if err != nil {
		gatewayErr := errors.NewInternalError(err)
		gatewayErr.WriteHTTP(w)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		gatewayErr := errors.NewInternalError(fmt.Errorf("internal service returned status: %d", resp.StatusCode))
		gatewayErr.WriteHTTP(w)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		gatewayErr := errors.NewInternalError(err)
		gatewayErr.WriteHTTP(w)
		return
	}

	var suggestions perfume.Suggestions
	if err := json.Unmarshal(body, &suggestions); err != nil {
		gatewayErr := errors.NewInternalError(err)
		gatewayErr.WriteHTTP(w)
		return
	}

	if len(suggestions.Perfumes) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		writeNoContentResponse(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		log.Printf("Error writing response: %v\n", err)
	}
}

func getTimeoutFromRequest(r http.Request) time.Duration {
	if strings.EqualFold(r.URL.Query().Get("use_ai"), "true") {
		timeout, err := time.ParseDuration(os.Getenv(aITimeout))
		if err != nil {
			return defaultAITimeout
		}
		return timeout
	}
	timeout, err := time.ParseDuration(os.Getenv(nonAITimeout))
	if err != nil {
		return defaultNonAITimeout
	}
	return timeout
}
func getHTTPClient(timeout time.Duration) *http.Client {
	responseHeaderTimeout := max(timeout-1*time.Second, 1*time.Second)

	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:          10,
			IdleConnTimeout:       30 * time.Second,
			DisableCompression:    true,
			ResponseHeaderTimeout: responseHeaderTimeout,
			DisableKeepAlives:     false,
		},
	}
}

func writeNoContentResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	response := map[string]string{
		"message": "No recommendations available",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding no content response: %v\n", err)
	}
}
