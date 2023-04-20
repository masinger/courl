package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func Read(path string) (config *Config, err error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	defer func() {
		if tempErr := f.Close(); tempErr != nil {
			err = tempErr
		}
	}()

	var result Config
	err = yaml.NewDecoder(f).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func ReadDefault() (*Config, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configLocation := filepath.Join(userHome, ".courl", "config")
	return Read(configLocation)
}
