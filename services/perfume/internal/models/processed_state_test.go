package models

import (
	"encoding/json"
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/errors"
)

func TestNewProcessedState(t *testing.T) {
	state := NewProcessedState()

	if !state.Success {
		t.Errorf("NewProcessedState().Success = false, want true")
	}
	if state.SuccessfulCount != 0 {
		t.Errorf("NewProcessedState().SuccessfulCount = %d, want 0", state.SuccessfulCount)
	}
	if state.FailedCount != 0 {
		t.Errorf("NewProcessedState().FailedCount = %d, want 0", state.FailedCount)
	}
	if state.Error != nil {
		t.Errorf("NewProcessedState().Error = %v, want nil", state.Error)
	}
}

func TestProcessedState_JSONSerialization(t *testing.T) {
	state := NewProcessedState()
	state.SuccessfulCount = 5
	state.FailedCount = 2

	jsonData, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("Failed to marshal ProcessedState: %v", err)
	}

	var unmarshaled ProcessedState
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ProcessedState: %v", err)
	}

	if unmarshaled.Success != state.Success {
		t.Errorf("Unmarshaled Success = %v, want %v", unmarshaled.Success, state.Success)
	}
	if unmarshaled.SuccessfulCount != state.SuccessfulCount {
		t.Errorf("Unmarshaled SuccessfulCount = %d, want %d", unmarshaled.SuccessfulCount, state.SuccessfulCount)
	}
	if unmarshaled.FailedCount != state.FailedCount {
		t.Errorf("Unmarshaled FailedCount = %d, want %d", unmarshaled.FailedCount, state.FailedCount)
	}
}

func TestProcessedState_ErrorFieldNotSerialized(t *testing.T) {
	state := ProcessedState{
		Success:         false,
		SuccessfulCount: 0,
		FailedCount:    0,
		Error:           errors.NewDBError("test error", nil),
	}

	jsonData, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("Failed to marshal ProcessedState: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Проверяем, что поле Error не присутствует в JSON
	if _, exists := result["error"]; exists {
		t.Error("Error field should not be serialized to JSON")
	}

	// Проверяем, что другие поля присутствуют
	if _, exists := result["success"]; !exists {
		t.Error("Success field should be present in JSON")
	}
	if _, exists := result["successful_count"]; !exists {
		t.Error("SuccessfulCount field should be present in JSON")
	}
	if _, exists := result["failed_count"]; !exists {
		t.Error("FailedCount field should be present in JSON")
	}
}

func TestProcessedState_WithError(t *testing.T) {
	dbErr := errors.NewDBError("connection failed", nil)
	state := ProcessedState{
		Success:         false,
		SuccessfulCount: 0,
		FailedCount:    0,
		Error:           dbErr,
	}

	if state.Error == nil {
		t.Error("Expected Error to be set")
	}

	if state.Error != dbErr {
		t.Errorf("Error = %v, want %v", state.Error, dbErr)
	}

	// Проверяем, что Error не влияет на JSON сериализацию
	jsonData, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("Failed to marshal ProcessedState with error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if _, exists := result["error"]; exists {
		t.Error("Error field should not be serialized to JSON even when set")
	}
}

