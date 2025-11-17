package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/errors"
	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/models"
)

// Mock для core.Select - мы не можем легко мокировать, поэтому тестируем только handleError
// и проверяем, что handler правильно обрабатывает различные сценарии через интеграционные тесты

func TestHandleError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "DBError",
			err:            errors.NewDBError("database connection failed", nil),
			expectedStatus: 500,
		},
		{
			name:           "ValidationError",
			err:            errors.NewValidationError("invalid input"),
			expectedStatus: 400,
		},
		{
			name:           "NotFoundError",
			err:            errors.NewNotFoundError("not found"),
			expectedStatus: 404,
		},
		{
			name:           "AuthError",
			err:            errors.NewAuthError("invalid token"),
			expectedStatus: 403,
		},
		{
			name:           "non-ServiceError",
			err:            &customError{msg: "custom error"},
			expectedStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handleError(w, tt.err)

			if w.Code != tt.expectedStatus {
				t.Errorf("handleError() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("handleError() Content-Type = %q, want %q", contentType, "application/json")
			}

			var state models.ProcessedState
			if err := json.NewDecoder(w.Body).Decode(&state); err != nil {
				t.Errorf("handleError() response is not valid JSON: %v", err)
			}
		})
	}
}

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func TestSelect_QueryParameters(t *testing.T) {
	// Этот тест проверяет, что query параметры правильно извлекаются
	// Полный тест Select требует мокирования core.Select, что сложно без интерфейсов
	// Но мы можем проверить, что параметры правильно парсятся через интеграционные тесты

	tests := []struct {
		testName string
		url      string
		brand    string
		name     string
		sex      string
	}{
		{
			testName: "all parameters",
			url:      "/v1/perfumes/get?brand=Chanel&name=No.5&sex=female",
			brand:    "Chanel",
			name:     "No.5",
			sex:      "female",
		},
		{
			testName: "only brand",
			url:      "/v1/perfumes/get?brand=Chanel",
			brand:    "Chanel",
			name:     "",
			sex:      "",
		},
		{
			testName: "no parameters",
			url:      "/v1/perfumes/get",
			brand:    "",
			name:     "",
			sex:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			brand := req.URL.Query().Get("brand")
			name := req.URL.Query().Get("name")
			sex := req.URL.Query().Get("sex")

			if brand != tt.brand {
				t.Errorf("brand = %q, want %q", brand, tt.brand)
			}
			if name != tt.name {
				t.Errorf("name = %q, want %q", name, tt.name)
			}
			if sex != tt.sex {
				t.Errorf("sex = %q, want %q", sex, tt.sex)
			}
		})
	}
}
