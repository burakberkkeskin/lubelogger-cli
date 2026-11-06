package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig_Precedence(t *testing.T) {
	// Clean ambient environment for the duration of this test
	origURL := os.Getenv("LUBELOGGER_URL")
	origKey := os.Getenv("LUBELOGGER_API_KEY")
	origTimeout := os.Getenv("LUBELOGGER_TIMEOUT")
	defer func() {
		_ = os.Setenv("LUBELOGGER_URL", origURL)
		_ = os.Setenv("LUBELOGGER_API_KEY", origKey)
		_ = os.Setenv("LUBELOGGER_TIMEOUT", origTimeout)
	}()

	_ = os.Unsetenv("LUBELOGGER_URL")
	_ = os.Unsetenv("LUBELOGGER_API_KEY")
	_ = os.Unsetenv("LUBELOGGER_TIMEOUT")

	// Test Default
	cfg := LoadConfig("", "", "", "", "", false)
	if cfg.URL != "http://localhost:5000" {
		t.Errorf("expected default URL http://localhost:5000, got %s", cfg.URL)
	}

	// Test Env override
	_ = os.Setenv("LUBELOGGER_URL", "http://env-host:8080")
	_ = os.Setenv("LUBELOGGER_API_KEY", "env-api-key")
	_ = os.Setenv("LUBELOGGER_TIMEOUT", "15s")

	cfgEnv := LoadConfig("", "", "", "", "", false)
	if cfgEnv.URL != "http://env-host:8080" {
		t.Errorf("expected env URL http://env-host:8080, got %s", cfgEnv.URL)
	}
	if cfgEnv.APIKey != "env-api-key" {
		t.Errorf("expected env APIKey env-api-key, got %s", cfgEnv.APIKey)
	}
	if cfgEnv.Timeout != 15*time.Second {
		t.Errorf("expected env Timeout 15s, got %v", cfgEnv.Timeout)
	}

	// Test Flag override (highest precedence)
	cfgFlag := LoadConfig("http://flag-host:9000", "flag-api-key", "", "", "5s", true)
	if cfgFlag.URL != "http://flag-host:9000" {
		t.Errorf("expected flag URL http://flag-host:9000, got %s", cfgFlag.URL)
	}
	if cfgFlag.APIKey != "flag-api-key" {
		t.Errorf("expected flag APIKey flag-api-key, got %s", cfgFlag.APIKey)
	}
	if cfgFlag.Timeout != 5*time.Second {
		t.Errorf("expected flag Timeout 5s, got %v", cfgFlag.Timeout)
	}
	if !cfgFlag.JSONOutput {
		t.Errorf("expected JSONOutput to be true when flag is passed")
	}
}
