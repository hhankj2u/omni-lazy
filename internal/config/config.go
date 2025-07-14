package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Ollama OllamaConfig `json:"ollama"`
}

// OllamaConfig represents Ollama-specific configuration
type OllamaConfig struct {
	URL                       string  `json:"url"`
	Model                     string  `json:"model"`
	SeedOrNegative            int     `json:"seed_or_negative"`
	TemperatureIfNegativeSeed float64 `json:"temperature_if_negative_seed"`
	PullTimeout               int     `json:"pull_timeout"`
	HTTPTimeout               int     `json:"http_timeout"`
	TrimSpace                 bool    `json:"trim_space"`
	Verbose                   bool    `json:"verbose"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Ollama: OllamaConfig{
			URL:                       "http://localhost:11434",
			Model:                     "gemma3:4b",
			SeedOrNegative:            0,
			TemperatureIfNegativeSeed: 0.0,
			PullTimeout:               0,
			HTTPTimeout:               0,
			TrimSpace:                 true,
			Verbose:                   false,
		},
	}
}

// LoadConfig loads configuration from file or creates default if not exists
func LoadConfig(configPath string) (*Config, error) {
	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Try to load existing config
	if _, err := os.Stat(configPath); err == nil {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		var config Config
		if err := json.NewDecoder(file).Decode(&config); err != nil {
			return nil, fmt.Errorf("failed to decode config file: %w", err)
		}

		return &config, nil
	}

	// Create default config if file doesn't exist
	config := DefaultConfig()
	if err := SaveConfig(configPath, config); err != nil {
		return nil, fmt.Errorf("failed to save default config: %w", err)
	}

	return config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(configPath string, config *Config) error {
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

// GetConfigPath returns the default config file path
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return "config.json"
	}
	return filepath.Join(homeDir, ".omni-lazy", "config.json")
}
