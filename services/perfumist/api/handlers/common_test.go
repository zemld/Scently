package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/zemld/Scently/models"
	"github.com/zemld/Scently/perfumist/internal/config"
	"github.com/zemld/Scently/perfumist/internal/errors"
)

func TestGeneralParseSimilarParameters_ValidParams(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/?brand=Chanel&name=No5&sex=female", nil)
	params, err := generalParseSimilarParameters(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if params.Brand != "Chanel" {
		t.Fatalf("expected brand %q, got %q", "Chanel", params.Brand)
	}
	if params.Name != "No5" {
		t.Fatalf("expected name %q, got %q", "No5", params.Name)
	}
	if params.Sex != "female" {
		t.Fatalf("expected sex %q, got %q", "female", params.Sex)
	}
}

func TestGeneralParseSimilarParameters_MaleSex(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/?brand=Dior&name=Sauvage&sex=male", nil)
	params, err := generalParseSimilarParameters(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if params.Sex != "male" {
		t.Fatalf("expected sex %q, got %q", "male", params.Sex)
	}
}

func TestGeneralParseSimilarParameters_InvalidSexDefaultsToEmpty(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/?brand=Tom+Ford&name=Black+Orchid&sex=invalid", nil)
	params, err := generalParseSimilarParameters(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// WithSex only sets sex if it's "male" or "female", otherwise it leaves it empty
	if params.Sex != "" {
		t.Fatalf("expected empty sex, got %q", params.Sex)
	}
}

func TestGeneralParseSimilarParameters_EmptySexDefaultsToEmpty(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/?brand=Chanel&name=No5", nil)
	params, err := generalParseSimilarParameters(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// WithSex only sets sex if it's "male" or "female", otherwise it leaves it empty
	if params.Sex != "" {
		t.Fatalf("expected empty sex, got %q", params.Sex)
	}
}

func TestGeneralParseSimilarParameters_MissingBrand(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/?name=No5", nil)
	_, err := generalParseSimilarParameters(req)

	if err == nil {
		t.Fatal("expected validation error for missing brand")
	}
	validationErr, ok := err.(*errors.ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if validationErr.Field != "brand" {
		t.Fatalf("expected field %q, got %q", "brand", validationErr.Field)
	}
}

func TestGeneralParseSimilarParameters_MissingName(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/?brand=Chanel", nil)
	_, err := generalParseSimilarParameters(req)

	if err == nil {
		t.Fatal("expected validation error for missing name")
	}
	validationErr, ok := err.(*errors.ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if validationErr.Field != "name" {
		t.Fatalf("expected field %q, got %q", "name", validationErr.Field)
	}
}

func TestGeneralParseSimilarParameters_EmptyParams(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := generalParseSimilarParameters(req)

	if err == nil {
		t.Fatal("expected validation error for empty params")
	}
}

func TestCreatePerfumeHubFetcher_Success(t *testing.T) {
	t.Parallel()

	mockCM := &config.MockConfigManager{
		GetStringFunc: func(key string) (string, error) {
			switch key {
			case "get_perfumes_url":
				return "http://test:8000/v1/perfumes", nil
			case "perfume_hub_internal_token_env_name":
				return "TEST_TOKEN", nil
			default:
				return "", nil
			}
		},
	}

	originalToken := os.Getenv("TEST_TOKEN")
	defer func() {
		if originalToken != "" {
			os.Setenv("TEST_TOKEN", originalToken)
		} else {
			os.Unsetenv("TEST_TOKEN")
		}
	}()
	os.Setenv("TEST_TOKEN", "test-token-value")

	fetcher, err := createPerfumeHubFetcher(mockCM)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fetcher == nil {
		t.Fatal("expected non-nil fetcher")
	}
}

func TestCreatePerfumeHubFetcher_GetPerfumesUrlError(t *testing.T) {
	t.Parallel()

	mockCM := &config.MockConfigManager{
		GetStringFunc: func(key string) (string, error) {
			if key == "get_perfumes_url" {
				return "", errors.NewServiceError("config error", nil)
			}
			return "", nil
		},
	}

	fetcher, err := createPerfumeHubFetcher(mockCM)

	if err == nil {
		t.Fatal("expected error when get_perfumes_url fails")
	}
	if fetcher != nil {
		t.Fatal("expected nil fetcher on error")
	}
	serviceErr, ok := err.(*errors.ServiceError)
	if !ok {
		t.Fatalf("expected ServiceError, got %T", err)
	}
	if serviceErr.Message != "failed to get get_perfumes_url" {
		t.Fatalf("expected message %q, got %q", "failed to get get_perfumes_url", serviceErr.Message)
	}
}

func TestCreatePerfumeHubFetcher_PerfumeHubInternalTokenEnvError(t *testing.T) {
	t.Parallel()

	mockCM := &config.MockConfigManager{
		GetStringFunc: func(key string) (string, error) {
			switch key {
			case "get_perfumes_url":
				return "http://test:8000/v1/perfumes", nil
			case "perfume_hub_internal_token_env_name":
				return "", errors.NewServiceError("config error", nil)
			default:
				return "", nil
			}
		},
	}

	fetcher, err := createPerfumeHubFetcher(mockCM)

	if err == nil {
		t.Fatal("expected error when perfume_hub_internal_token_env_name fails")
	}
	if fetcher != nil {
		t.Fatal("expected nil fetcher on error")
	}
	serviceErr, ok := err.(*errors.ServiceError)
	if !ok {
		t.Fatalf("expected ServiceError, got %T", err)
	}
	if serviceErr.Message != "failed to get perfume_hub_internal_token_env_name" {
		t.Fatalf("expected message %q, got %q", "failed to get perfume_hub_internal_token_env_name", serviceErr.Message)
	}
}

func TestCreateAIFetcher_Success(t *testing.T) {
	t.Parallel()

	originalBaseURL := os.Getenv("BASE_URL")
	originalFolderID := os.Getenv("FOLDER_ID")
	originalModelName := os.Getenv("MODEL_NAME")
	originalAPIKey := os.Getenv("API_KEY")
	defer func() {
		if originalBaseURL != "" {
			os.Setenv("BASE_URL", originalBaseURL)
		} else {
			os.Unsetenv("BASE_URL")
		}
		if originalFolderID != "" {
			os.Setenv("FOLDER_ID", originalFolderID)
		} else {
			os.Unsetenv("FOLDER_ID")
		}
		if originalModelName != "" {
			os.Setenv("MODEL_NAME", originalModelName)
		} else {
			os.Unsetenv("MODEL_NAME")
		}
		if originalAPIKey != "" {
			os.Setenv("API_KEY", originalAPIKey)
		} else {
			os.Unsetenv("API_KEY")
		}
	}()

	os.Setenv("BASE_URL", "http://test:8000")
	os.Setenv("FOLDER_ID", "test-folder")
	os.Setenv("MODEL_NAME", "test-model")
	os.Setenv("API_KEY", "test-key")

	mockCM := &config.MockConfigManager{}
	fetcher := createAIFetcher(mockCM)

	if fetcher == nil {
		t.Fatal("expected non-nil fetcher")
	}
}

func TestHandleError_ValidationError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	validationErr := errors.NewValidationError("brand", "is required")
	handleError(w, validationErr)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if response.Error == "" {
		t.Fatal("expected non-empty error message")
	}
	if response.Error != validationErr.Error() {
		t.Fatalf("expected error message %q, got %q", validationErr.Error(), response.Error)
	}
}

func TestHandleError_NotFoundError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	notFoundErr := errors.NewNotFoundError("perfume not found")
	handleError(w, notFoundErr)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if response.Error != notFoundErr.Error() {
		t.Fatalf("expected error message %q, got %q", notFoundErr.Error(), response.Error)
	}
}

func TestHandleError_ServiceError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	serviceErr := errors.NewServiceError("internal service error", nil)
	handleError(w, serviceErr)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if response.Error != serviceErr.Error() {
		t.Fatalf("expected error message %q, got %q", serviceErr.Error(), response.Error)
	}
}

func TestHandleError_UnknownError(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	unknownErr := &customError{msg: "unknown error"}
	handleError(w, unknownErr)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if response.Error != "internal server error" {
		t.Fatalf("expected error message %q, got %q", "internal server error", response.Error)
	}
}

func TestWriteResponse_Success(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	response := SuggestResponse{
		Suggested: []models.Ranked{
			{
				Perfume: models.Perfume{Brand: "Chanel", Name: "No5", Sex: "female"},
				Score:   0.95,
				Rank:    1,
			},
		},
	}
	WriteResponse(w, response, http.StatusOK)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected Content-Type %q, got %q", "application/json", w.Header().Get("Content-Type"))
	}

	var decodedResponse SuggestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &decodedResponse); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(decodedResponse.Suggested) != 1 {
		t.Fatalf("expected 1 suggested perfume, got %d", len(decodedResponse.Suggested))
	}
	if decodedResponse.Suggested[0].Perfume.Brand != "Chanel" {
		t.Fatalf("expected brand %q, got %q", "Chanel", decodedResponse.Suggested[0].Perfume.Brand)
	}
}

func TestWriteResponse_ErrorResponse(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	response := ErrorResponse{Error: "test error"}
	WriteResponse(w, response, http.StatusBadRequest)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected Content-Type %q, got %q", "application/json", w.Header().Get("Content-Type"))
	}

	var decodedResponse ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &decodedResponse); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if decodedResponse.Error != "test error" {
		t.Fatalf("expected error %q, got %q", "test error", decodedResponse.Error)
	}
}

func TestWriteResponse_EmptyResponse(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	WriteResponse(w, nil, http.StatusNoContent)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected Content-Type %q, got %q", "application/json", w.Header().Get("Content-Type"))
	}
}

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}
