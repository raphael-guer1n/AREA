package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadAllServiceConfigs(rootDir string) ([]ServiceConfig, error) {
	var services []ServiceConfig

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "service.config.json" {
			serviceCfg, err := loadSingleConfig(path)
			if err != nil {
				return fmt.Errorf("error loading %s: %w", path, err)
			}
			services = append(services, serviceCfg)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return services, nil
}

func loadSingleConfig(path string) (ServiceConfig, error) {
	var cfg ServiceConfig

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(fileBytes, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return cfg, nil
}
