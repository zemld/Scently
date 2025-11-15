package handlers

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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
		log.Printf("Error creating request: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.URL.RawQuery = r.URL.Query().Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		} else if ctx.Err() == context.Canceled {
			http.Error(w, "Request canceled", http.StatusRequestTimeout)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	log.Printf("Response status: %s\n", resp.Status)
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
	log.Printf("Response body: %s\n", resp.Body)
}

func getTimeoutFromRequest(r http.Request) time.Duration {
	if r.URL.Query().Get("use_ai") == "true" {
		timeout, err := time.ParseDuration(aITimeout)
		if err != nil {
			return defaultAITimeout
		}
		return timeout
	}
	timeout, err := time.ParseDuration(nonAITimeout)
	if err != nil {
		return defaultNonAITimeout
	}
	return timeout
}
