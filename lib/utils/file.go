package utils

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"strings"
)

func GetAppDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".envme")
	err = EnsureDir(dir)
	if err != nil {
		return dir, err
	}
	return dir, nil
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

	file := filepath.Join(appDir, "config.yaml")
	if err := EnsureFile(file); err != nil {
		return "", err
	}

	return file, nil
}

func GetServiceDir(name string) (string, error) {
	appDir, err := GetAppDir()
	if err != nil {
		return "", err
	}

	srvDir := filepath.Join(appDir, name)
	if err := EnsureDir(srvDir); err != nil {
		return "", err
	}
	return srvDir, nil
}

func WriteComposeFile(name string, content []byte) error {
	dir, err := GetServiceDir(name)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "docker-compose.yaml"), content, 0644)
}

func WriteDockerfile(dir string, template string) error {
	fmt.Printf("Writing Dockerfile template %s to %s", template, dir)
	// TODO: Add template logic
	return os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte(template), 0644)
}

func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func EnsureFile(file string) error {
	if err := EnsureDir(filepath.Dir(file)); err != nil {
		return err
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return os.WriteFile(file, []byte(""), 0644)
	}
	return nil
}

func GetAbsPath(dir string) (string, error) {
	if strings.HasPrefix(dir, "..") {
		current, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return filepath.Join(current, dir), nil
	}

	if strings.HasPrefix(dir, ".") {
		current, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return strings.Replace(dir, ".", current, 1), nil
	}

	if strings.HasPrefix(dir, "~") {
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		return strings.Replace(dir, "~", home, 1), nil
	}

	return dir, nil
}
