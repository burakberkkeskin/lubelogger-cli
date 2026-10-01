package config

import (
	"os"
	"strings"
	"time"
)

// Config holds runtime configuration options.
type Config struct {
	URL        string
	APIKey     string
	Username   string
	Password   string
	Timeout    time.Duration
	JSONOutput bool
}

// IsTerminal returns true if stdout is an interactive terminal.
func IsTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// LoadConfig merges defaults, environment variables, and flag overrides.
func LoadConfig(flagURL string, flagAPIKey string, flagUsername string, flagPassword string, flagTimeout string, flagJSON bool) *Config {
	// 1. Defaults
	cfg := &Config{
		URL:        "http://localhost:5000",
		Timeout:    10 * time.Second,
		JSONOutput: !IsTerminal(), // If piped or captured by an agent, default to JSON
	}

	// 2. Environment Variables
	if envURL := os.Getenv("LUBELOGGER_URL"); envURL != "" {
		cfg.URL = envURL
	}
	if envAPIKey := os.Getenv("LUBELOGGER_API_KEY"); envAPIKey != "" {
		cfg.APIKey = envAPIKey
	}
	if envUser := os.Getenv("LUBELOGGER_USERNAME"); envUser != "" {
		cfg.Username = envUser
	}
	if envPass := os.Getenv("LUBELOGGER_PASSWORD"); envPass != "" {
		cfg.Password = envPass
	}
	if envFormat := strings.ToLower(os.Getenv("LUBELOGGER_OUTPUT_FORMAT")); envFormat != "" {
		if envFormat == "json" {
			cfg.JSONOutput = true
		} else if envFormat == "table" || envFormat == "text" {
			cfg.JSONOutput = false
		}
	}
	if envTimeout := os.Getenv("LUBELOGGER_TIMEOUT"); envTimeout != "" {
		if d, err := time.ParseDuration(envTimeout); err == nil {
			cfg.Timeout = d
		}
	}

	// 3. Command Line Flags (highest precedence)
	if flagURL != "" {
		cfg.URL = flagURL
	}
	if flagAPIKey != "" {
		cfg.APIKey = flagAPIKey
	}
	if flagUsername != "" {
		cfg.Username = flagUsername
	}
	if flagPassword != "" {
		cfg.Password = flagPassword
	}
	if flagTimeout != "" {
		if d, err := time.ParseDuration(flagTimeout); err == nil {
			cfg.Timeout = d
		}
	}
	if flagJSON {
		cfg.JSONOutput = true
	}

	return cfg
}
