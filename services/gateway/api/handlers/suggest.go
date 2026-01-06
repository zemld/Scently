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

	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/config"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/errors"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/models/perfume"
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, perfumistUrl, nil)
	if err != nil {
		gatewayErr := errors.NewInternalError(err)
		gatewayErr.WriteHTTP(w)
		return
	}
	req.URL.RawQuery = r.URL.Query().Encode()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("PERFUMIST_INTERNAL_TOKEN")))

	timeout := getTimeoutFromRequest(*r, m)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		gatewayErr := errors.NewInternalError(err)
		gatewayErr.WriteHTTP(w)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		gatewayErr := errors.NewInternalError(err)
		gatewayErr.WriteHTTP(w)
		return
	}

	switch resp.StatusCode {
	case http.StatusOK:
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

	case http.StatusBadRequest, http.StatusNotFound:
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			gatewayErr := errors.ErrBadRequest(fmt.Errorf("perfumist service error"))
			gatewayErr.WriteHTTP(w)
			return
		}
		gatewayErr := errors.ErrBadRequest(fmt.Errorf("%s", errorResponse.Error))
		gatewayErr.WriteHTTP(w)

	case http.StatusInternalServerError:
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			gatewayErr := errors.NewInternalError(fmt.Errorf("perfumist service returned 500"))
			gatewayErr.WriteHTTP(w)
			return
		}
		gatewayErr := errors.NewInternalError(fmt.Errorf("perfumist service error: %s", errorResponse.Error))
		gatewayErr.WriteHTTP(w)

	default:
		gatewayErr := errors.NewInternalError(fmt.Errorf("internal service returned status: %d", resp.StatusCode))
		gatewayErr.WriteHTTP(w)
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

func validateParameters(r http.Request) *errors.GatewayError {
	if r.URL.Query().Get("brand") == "" {
		return errors.ErrBadRequest(fmt.Errorf("brand is required"))
	}
	if r.URL.Query().Get("name") == "" {
		return errors.ErrBadRequest(fmt.Errorf("name is required"))
	}
	return nil
}
