package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// GetConfigPath - get configuration file path
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: User home directory not found:", err)
		return ""
	}

	return filepath.Join(homeDir, ".config", "aurora", "config.yaml")
}

// LoadConfig - load configuration file or create a new one
func LoadConfig() error {
	configPath := GetConfigPath()

	// Check and create configuration directory
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create configuration directory: %w", err)
		}
	}

	// Read configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, save default configuration
			CurrentConfig = DefaultConfig
			return SaveConfig()
		}
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Read YAML format
	if err := yaml.Unmarshal(data, &CurrentConfig); err != nil {
		return fmt.Errorf("failed to parse configuration file: %w", err)
	}

	return nil
}

// SaveConfig - save current configuration to file
func SaveConfig() error {
	configPath := GetConfigPath()

	// Convert to YAML format
	data, err := yaml.Marshal(&CurrentConfig)
	if err != nil {
		return fmt.Errorf("failed to convert configuration to YAML: %w", err)
	}

	// Save to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to save configuration file: %w", err)
	}

	return nil
}
