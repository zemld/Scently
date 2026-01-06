package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/errors"
	"github.com/zemld/PerfumeRecommendationSystem/gateway/internal/models/perfume"
)

func proxyRequestToPerfumist(
	ctx context.Context,
	perfumistUrl string,
	originalReq *http.Request,
	timeout time.Duration,
	requireAuth bool,
) (*http.Response, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, perfumistUrl, nil)
	if err != nil {
		return nil, nil, err
	}
	req.URL.RawQuery = originalReq.URL.Query().Encode()
	if requireAuth {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("PERFUMIST_INTERNAL_TOKEN")))
	}

	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, nil, err
	}

	return resp, body, nil
}

func handlePerfumistResponse(w http.ResponseWriter, resp *http.Response, body []byte) error {
	switch resp.StatusCode {
	case http.StatusOK:
		var suggestions perfume.Suggestions
		if err := json.Unmarshal(body, &suggestions); err != nil {
			gatewayErr := errors.NewInternalError(err)
			gatewayErr.WriteHTTP(w)
			return err
		}

		if len(suggestions.Perfumes) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			writeNoContentResponse(w)
			return nil
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(body); err != nil {
			return err
		}

	case http.StatusBadRequest, http.StatusNotFound:
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			gatewayErr := errors.ErrBadRequest(fmt.Errorf("perfumist service error"))
			gatewayErr.WriteHTTP(w)
			return err
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
			return err
		}
		gatewayErr := errors.NewInternalError(fmt.Errorf("perfumist service error: %s", errorResponse.Error))
		gatewayErr.WriteHTTP(w)

	default:
		gatewayErr := errors.NewInternalError(fmt.Errorf("internal service returned status: %d", resp.StatusCode))
		gatewayErr.WriteHTTP(w)
	}

	return nil
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
