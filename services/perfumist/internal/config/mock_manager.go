package config

import (
	"context"
	"time"

	"github.com/zemld/config-manager/pkg/cm"
)

type MockConfigManager struct {
	GetStringFunc              func(key string) (string, error)
	GetStringWithDefaultFunc   func(key string, defaultValue string) string
	GetIntFunc                 func(key string) (int, error)
	GetIntWithDefaultFunc      func(key string, defaultValue int) int
	GetFloatFunc               func(key string) (float64, error)
	GetFloatWithDefaultFunc    func(key string, defaultValue float64) float64
	GetDurationFunc            func(key string) (time.Duration, error)
	GetDurationWithDefaultFunc func(key string, defaultValue time.Duration) time.Duration
	GetBoolFunc                func(key string) (bool, error)
	GetBoolWithDefaultFunc     func(key string, defaultValue bool) bool
	LoadConfigFunc             func(ctx context.Context) error
	StartLoadingFunc           func(interval time.Duration)
	StopLoadingFunc            func()
}

func (m *MockConfigManager) GetString(key string) (string, error) {
	if m.GetStringFunc != nil {
		return m.GetStringFunc(key)
	}
	switch key {
	case "get_perfumes_url":
		return "http://test:8000/v1/perfumes", nil
	case "perfume_hub_internal_token_env_name":
		return "TEST_TOKEN", nil
	default:
		return "", nil
	}
}

func (m *MockConfigManager) GetStringWithDefault(key string, defaultValue string) string {
	if m.GetStringWithDefaultFunc != nil {
		return m.GetStringWithDefaultFunc(key, defaultValue)
	}
	return defaultValue
}

func (m *MockConfigManager) GetInt(key string) (int, error) {
	if m.GetIntFunc != nil {
		return m.GetIntFunc(key)
	}
	return 0, nil
}

func (m *MockConfigManager) GetIntWithDefault(key string, defaultValue int) int {
	if m.GetIntWithDefaultFunc != nil {
		return m.GetIntWithDefaultFunc(key, defaultValue)
	}
	return defaultValue
}

func (m *MockConfigManager) GetFloat(key string) (float64, error) {
	if m.GetFloatFunc != nil {
		return m.GetFloatFunc(key)
	}
	return 0.0, nil
}

func (m *MockConfigManager) GetFloatWithDefault(key string, defaultValue float64) float64 {
	if m.GetFloatWithDefaultFunc != nil {
		return m.GetFloatWithDefaultFunc(key, defaultValue)
	}
	return defaultValue
}

func (m *MockConfigManager) GetDuration(key string) (time.Duration, error) {
	if m.GetDurationFunc != nil {
		return m.GetDurationFunc(key)
	}
	return 0, nil
}

func (m *MockConfigManager) GetDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	if m.GetDurationWithDefaultFunc != nil {
		return m.GetDurationWithDefaultFunc(key, defaultValue)
	}
	return defaultValue
}

func (m *MockConfigManager) GetBool(key string) (bool, error) {
	if m.GetBoolFunc != nil {
		return m.GetBoolFunc(key)
	}
	return false, nil
}

func (m *MockConfigManager) GetBoolWithDefault(key string, defaultValue bool) bool {
	if m.GetBoolWithDefaultFunc != nil {
		return m.GetBoolWithDefaultFunc(key, defaultValue)
	}
	return defaultValue
}

func (m *MockConfigManager) LoadConfig(ctx context.Context) error {
	return nil
}

func (m *MockConfigManager) StartLoading(interval time.Duration) {}
func (m *MockConfigManager) StopLoading()                        {}

// Убедимся, что MockConfigManager реализует интерфейс cm.ConfigManager
var _ cm.ConfigManager = (*MockConfigManager)(nil)
