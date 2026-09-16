package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(path, configFileName)
	return fullPath, nil
}

func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, err
	}
	var newConfig Config
	err = json.Unmarshal(data, &newConfig)
	if err != nil {
		return Config{}, err
	}
	return newConfig, nil
}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	jsonCfg, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(fullPath, jsonCfg, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}
