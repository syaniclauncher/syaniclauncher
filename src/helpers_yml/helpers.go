package helpers_yml

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

func Read(path string, destination any) (bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	if err := yaml.Unmarshal(data, destination); err != nil {
		return true, err
	}
	return true, nil
}

func Marshal(value any) ([]byte, error) {
	return yaml.Marshal(value)
}

func Write(path string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(path, data, perm); err != nil {
		return err
	}
	return os.Chmod(path, perm)
}

func Quarantine(path string) (string, error) {
	backupPath := fmt.Sprintf("%s.%d.corrupt", path, time.Now().UnixNano())
	if err := os.Rename(path, backupPath); err != nil {
		return "", err
	}
	return backupPath, nil
}
