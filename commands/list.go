package commands

import (
	"os"
	"path/filepath"
)

func HandleList() ([]string, error) {
	baseDir := filepath.Join(os.Getenv("HOME"), ".config", "custom-ime")

	var projects []string

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return projects, nil // دایرکتوری وجود نداره = هیچ پروژه‌ای نیست
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			projects = append(projects, entry.Name())
		}
	}

	return projects, nil
}

func ProjectExists(projectName string) (bool, error) {
	projects, err := HandleList()
	if err != nil {
		return false, err
	}

	for _, project := range projects {
		if project == projectName {
			return true, nil
		}
	}
	return false, nil
}