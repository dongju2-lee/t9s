package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	TerraformRoot   string         `yaml:"terraform_root"`
	Backend         BackendConfig  `yaml:"backend"`
	Defaults        DefaultsConfig `yaml:"defaults"`
}

// BackendConfig represents the Terraform backend configuration
type BackendConfig struct {
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
	Prefix string `yaml:"prefix,omitempty"`
}

// DefaultsConfig represents default application settings
type DefaultsConfig struct {
	AutoRefresh     bool `yaml:"auto_refresh"`
	RefreshInterval int  `yaml:"refresh_interval"` // in seconds
}

// Load loads configuration from file
func Load() (*Config, error) {
	configPath := getConfigPath()

	// If config doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return createDefaultConfig(configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// Save saves configuration to file
func (c *Config) Save() error {
	configPath := getConfigPath()

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// getConfigPath returns the path to the config file
func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".p9s/config.yaml"
	}
	return filepath.Join(home, ".p9s", "config.yaml")
}

// createDefaultConfig creates a default configuration
func createDefaultConfig(path string) (*Config, error) {
	// Try to detect current directory as terraform root
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}

	config := &Config{
		TerraformRoot: currentDir,
		Backend: BackendConfig{
			Bucket: "terraform-state",
			Region: "ap-northeast-2",
		},
		Defaults: DefaultsConfig{
			AutoRefresh:     true,
			RefreshInterval: 60,
		},
	}

	// Save default config
	if err := config.Save(); err != nil {
		return nil, err
	}

	return config, nil
}
