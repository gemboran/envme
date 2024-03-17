package utils

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
)

func GetAppDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".envme"), nil
}

func GetListServices() ([]string, error) {
	appDir, err := GetAppDir()
	if err != nil {
		return nil, err
	}

	return filepath.Glob(filepath.Join(appDir, "*"))
}

func GetConfigFile() (string, error) {
	appDir, err := GetAppDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(appDir, "config.yaml"), nil
}

func GetServiceDir(name string) (string, error) {
	appDir, err := GetAppDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(appDir, name), nil
}

func WriteComposeFile(name string, content []byte) error {
	dir, err := GetServiceDir(name)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "docker-compose.yaml"), []byte(content), 0644)
}
