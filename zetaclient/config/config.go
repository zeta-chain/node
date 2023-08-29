package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const filename string = "zetaclient_config.json"
const folder string = "config"

func Save(config *Config, path string) error {
	folderPath := filepath.Join(path, folder)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}
	file := filepath.Join(path, folder, filename)
	file = filepath.Clean(file)

	jsonFile, _ := json.MarshalIndent(config, "", "    ")
	err = os.WriteFile(file, jsonFile, 0600)
	if err != nil {
		return err
	}
	return nil
}

func Load(path string) (*Config, error) {
	file := filepath.Join(path, folder, filename)
	file, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	file = filepath.Clean(file)
	cfg := NewConfig()
	input, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(input, &cfg)
	if err != nil {
		return nil, err
	}
	cfg.TssPath = GetPath(cfg.TssPath)
	cfg.PreParamsPath = GetPath(cfg.PreParamsPath)
	return cfg, nil
}

func GetPath(inputPath string) string {
	path := strings.Split(inputPath, "/")
	if len(path) > 0 {
		if path[0] == "~" {
			home, _ := os.UserHomeDir()
			path[0] = home
		}
	}
	return filepath.Join(path...)
}
