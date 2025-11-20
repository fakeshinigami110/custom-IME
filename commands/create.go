package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"ime-tool/utils"
)
func getTemplatePath() (string, error) {
    if _, err := os.Stat("templates/cmake_main.txt"); err == nil {
        return "templates", nil
    }
    
    possiblePaths := []string{
        "/usr/share/custom-ime/templates/",
        "/usr/local/share/custom-ime/templates/",
        "/opt/custom-ime/templates/",
    }
    
    for _, path := range possiblePaths {
        if _, err := os.Stat(filepath.Join(path, "cmake_main.txt")); err == nil {
            return path, nil
        }
    }
    
    return "", fmt.Errorf("cannot find templates directory")
}

func HandleCreate(cfg Config, forceOverwrite bool) error {

	
	fmt.Println("Creating new IME project...")
	db , err:= LoadDB()
	exists, err := ProjectExists(cfg.ProjectName)
	_ , exists2 := db.GetProject(cfg.ProjectName)
	if err != nil {
		return fmt.Errorf("error : %v" , err)
	}
	if !forceOverwrite {
		

		if exists {
			return fmt.Errorf("project '%s' already exists", cfg.ProjectName)
		}
		

		if exists2 {
			return fmt.Errorf("project '%s' already exists", cfg.ProjectName)
		}


	} else {
		if exists {
			projectDir := filepath.Join(os.Getenv("HOME"), ".config", "custom-ime", cfg.ProjectName)
			if err := os.RemoveAll(projectDir); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove existing project: %v", err)
			}
			fmt.Printf("%v Removed existing project: %s\n", ResultColor(true) ,cfg.ProjectName)
		}
		if exists2 {
			db.DeleteProject(cfg.ProjectName)
		}

	}

	if err := createIME(cfg); err != nil {
		return fmt.Errorf("failed to create IME: %v", err)
	}

	return nil
}

func createIME(cfg Config) error {
	fmt.Printf("[-] Creating project: %s\n", cfg.ProjectName)

	baseDir := filepath.Join(os.Getenv("HOME"), ".config", "custom-ime")
	projectDir := filepath.Join(baseDir, cfg.ProjectName)

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("failed to create base directory %s: %v", baseDir, err)
	}
	fmt.Printf("[Note] Projects base directory: %s\n", baseDir)

	dirs := []string{
		projectDir,
		filepath.Join(projectDir, "src"),
		filepath.Join(projectDir, "config"),
		filepath.Join(projectDir, "build"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
		fmt.Printf("%v Created directory: %s\n",ResultColor(true), dir)
	}

	templates := []struct {
		templateName string
		outputPath   string
	}{
		{"cmake_main.txt", filepath.Join(projectDir, "CMakeLists.txt")},
		{"cmake_src.txt", filepath.Join(projectDir, "src", "CMakeLists.txt")},
		{"ime.h", filepath.Join(projectDir, "src", cfg.IMEName+".h")},
		{"ime.cpp", filepath.Join(projectDir, "src", cfg.IMEName+".cpp")},
		{"addon.conf", filepath.Join(projectDir, "src", cfg.IMEName+"-addon.conf.in")},
		{"ime.conf", filepath.Join(projectDir, "src", cfg.IMEName+".conf")},
	}

	path , err :=getTemplatePath()
	if err != nil {
		return err
	}

	for _, tmpl := range templates {
		tmplPath := filepath.Join(path, tmpl.templateName)
		if err := generateFromTemplate(tmplPath, tmpl.outputPath, cfg); err != nil {
			return fmt.Errorf("failed to generate %s: %v", tmpl.templateName, err)
		}
		fmt.Printf("%v Generated: %s\n", ResultColor(true),tmpl.outputPath)
	}

	configDest := filepath.Join(projectDir, "config", cfg.IMEName+".conf")
	if cfg.ConfigFile != "" {
		if err := utils.CopyFile(cfg.ConfigFile, configDest); err != nil {
			return fmt.Errorf("failed to copy config file: %v", err)
		}
		fmt.Printf("%v Copied config: %s → %s\n",ResultColor(true), cfg.ConfigFile, configDest)
	} else {
		if err := createDefaultConfig(cfg, configDest); err != nil {
			return fmt.Errorf("failed to create default config: %v", err)
		}
		fmt.Printf("%v Created default config: %s\n",ResultColor(true), configDest)
	}
	db , err:= LoadDB()
	if err != nil {
		return fmt.Errorf("an error ocured during loading IME database : %v" , err)
	}
	cfg.Installed = false
	cfg.ProjectPath = projectDir
	
	err = db.AddProject(cfg)
	if err != nil {
		return fmt.Errorf("an error ocured during adding IME to database : %v" , err)
	}

	fmt.Printf("\n\n%v IME project '%s' created successfully!\n", ResultColor(true),cfg.ProjectName)
	fmt.Printf("%v Project location: %s\n",ResultColor(true), projectDir)
	fmt.Printf("%v All projects are stored in: %s\n", ResultColor(true),baseDir)
	fmt.Printf("\n%v\n" , ColorIt("Next steps:" , "blue"))
	fmt.Printf("  custom-ime install -p %s\n", cfg.ProjectName)
	
	return nil
}

func generateFromTemplate(templatePath, outputPath string, cfg Config) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %v", templatePath, err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", outputPath, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, cfg); err != nil {
		return fmt.Errorf("failed to execute template %s: %v", templatePath, err)
	}

	return nil
}

func createDefaultConfig(cfg Config, destPath string) error {
	defaultConfig := fmt.Sprintf(`# %s Configuration
# Auto-generated by omikami sama

[Settings]
# true / false
convert_numbers_to_binary=false

# keep / ignore
unknown_chars_behavior=keep
add_spaces=true

# leave it blank for no number separator
number_separator=$

# true / false
case_sensitive=false

[Characters]

[Capitals]

[Digits]

[Keywords]

[Operators]


`, cfg.IMEName)

	return os.WriteFile(destPath, []byte(defaultConfig), 0644)
}