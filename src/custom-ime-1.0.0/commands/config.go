// config.go
package commands

import (
	"strings"
	"path/filepath"
	"os"
	"time"
)

type Config struct {
    ProjectName string    `json:"name"`
    IMEName     string    `json:"ime_name"`
    Label       string    `json:"label"`
    Icon        string    `json:"icon"`
    LangCode    string    `json:"lang_code"`
    Description string    `json:"description"`
    ConfigFile  string    `json:"config_file,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Installed   bool      `json:"installed"`
    InstallPath string    `json:"install_path,omitempty"`
    ProjectPath string    `json:"project_path,omitempty"`
}


type Projects map[string]Config


func (c Config) ProjectNameUpper() string {
	return strings.ToUpper(c.ProjectName)
}

func (c Config) IMENameUpper() string {
	return strings.ToUpper(c.IMEName)
}

func (c Config) GetFullProjectPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "custom-ime", c.ProjectName)
}