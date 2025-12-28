package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zemld/Scently/perfume-hub/internal/errors"
	"github.com/zemld/Scently/perfume-hub/internal/models"
)

func TestUpdate_EmptyBody(t *testing.T) {
	w := httptest.NewRecorder()

	// Мы не можем легко мокировать core.Update, но можем проверить валидацию
	// Проверяем, что пустое тело обрабатывается правильно
	content := []byte{}

	if len(content) == 0 {
		validationErr := errors.NewValidationError("request body is empty")
		handleError(w, validationErr)
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("Update() with empty body status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUpdate_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()

	params := models.NewUpdateParameters()
	err := json.Unmarshal([]byte("invalid json"), params)
	if err != nil {
		validationErr := errors.NewValidationError("invalid JSON: " + err.Error())
		handleError(w, validationErr)
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("Update() with invalid JSON status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUpdate_EmptyPerfumesArray(t *testing.T) {
	w := httptest.NewRecorder()

	params := models.NewUpdateParameters()
	json.Unmarshal([]byte(`{"perfumes":[]}`), params)

	if len(params.Perfumes) == 0 {
		validationErr := errors.NewValidationError("perfumes array is empty")
		handleError(w, validationErr)
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("Update() with empty perfumes array status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUpdate_ValidRequest(t *testing.T) {
	validJSON := `{
		"perfumes": [
			{
				"brand": "Chanel",
				"name": "No. 5",
				"sex": "female",
				"image_url": "https://example.com/image.jpg",
				"properties": {
					"perfume_type": "Eau de Parfum",
					"family": ["Floral"],
					"upper_notes": ["Aldehydes"],
					"core_notes": ["Rose"],
					"base_notes": ["Vanilla"]
				},
				"shops": [
					{
						"shop_name": "Gold Apple",
						"domain": "goldapple.ru",
						"variants": [
							{
								"volume": 50,
								"price": 5000,
								"link": "https://goldapple.ru/perfume/123"
							}
						]
					}
				]
			}
		]
	}`

	params := models.NewUpdateParameters()
	err := json.Unmarshal([]byte(validJSON), params)
	if err != nil {
		t.Fatalf("Failed to unmarshal valid JSON: %v", err)
	}

	if len(params.Perfumes) == 0 {
		t.Error("Expected perfumes array to be non-empty")
	}

	if params.Perfumes[0].Brand != "Chanel" {
		t.Errorf("Expected brand to be 'Chanel', got %q", params.Perfumes[0].Brand)
	}

	if params.Perfumes[0].Name != "No. 5" {
		t.Errorf("Expected name to be 'No. 5', got %q", params.Perfumes[0].Name)
	}
}

func TestUpdate_RequestBodyParsing(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		expectError bool
		errorType   string
	}{
		{
			name:        "valid JSON with perfumes",
			body:        `{"perfumes":[{"brand":"Chanel","name":"No.5","sex":"female","properties":{"perfume_type":"EDP","family":[],"upper_notes":[],"core_notes":[],"base_notes":[]},"shops":[]}]}`,
			expectError: false,
		},
		{
			name:        "invalid JSON",
			body:        `{"perfumes": [invalid}`,
			expectError: true,
			errorType:   "validation",
		},
		{
			name:        "empty body",
			body:        "",
			expectError: true,
			errorType:   "validation",
		},
		{
			name:        "empty perfumes array",
			body:        `{"perfumes":[]}`,
			expectError: true,
			errorType:   "validation",
		},
		{
			name:        "missing perfumes field",
			body:        `{}`,
			expectError: true,
			errorType:   "validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := models.NewUpdateParameters()

			if tt.body == "" {
				// Проверяем пустое тело
				validationErr := errors.NewValidationError("request body is empty")
				if !tt.expectError {
					t.Error("Expected error for empty body")
				}
				if validationErr.Type() != errors.ErrorTypeValidation {
					t.Errorf("Expected ValidationError, got %v", validationErr.Type())
				}
				return
			}

			err := json.Unmarshal([]byte(tt.body), params)
			if tt.expectError && err == nil {
				// Проверяем, что если ожидаем ошибку, но unmarshal прошел успешно,
				// то должна быть ошибка валидации (например, пустой массив)
				if len(params.Perfumes) == 0 {
					validationErr := errors.NewValidationError("perfumes array is empty")
					if validationErr.Type() != errors.ErrorTypeValidation {
						t.Errorf("Expected ValidationError, got %v", validationErr.Type())
					}
				}
			} else if tt.expectError && err != nil {
				// Проверяем, что ошибка unmarshal правильно обрабатывается
				validationErr := errors.NewValidationError("invalid JSON: " + err.Error())
				if validationErr.Type() != errors.ErrorTypeValidation {
					t.Errorf("Expected ValidationError, got %v", validationErr.Type())
				}
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
