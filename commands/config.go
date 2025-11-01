package commands

import (
	"strings"
	"path/filepath"
	"os"
)

// Config ساختار مشترک برای همه کامندها
type Config struct {
	ProjectName string
	IMEName     string
	Label       string
	Icon        string
	LangCode    string
	Description string
	ConfigFile  string
}

func (c Config) ProjectNameUpper() string {
	return strings.ToUpper(c.ProjectName)
}

func (c Config) IMENameUpper() string {
	return strings.ToUpper(c.IMEName)
}

func (c Config) GetFullProjectPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "custom-ime", c.ProjectName)
}