package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	ProjectName string
	IMEName     string
	Label       string
	Icon        string
	LangCode    string
	Description string
	ConfigFile  string
	ProjectDir  string
}

func handleDualFlags(flagShort, flagLong *string, defaultValue string) *string {
	if *flagShort == defaultValue && *flagLong == defaultValue {
		return flagShort 
	}
	
	if *flagShort != defaultValue && *flagLong != defaultValue && *flagShort != *flagLong {
		fmt.Println("Error: You can't pass two different values for the same flag")
		printUsage()
		os.Exit(1)
	}
	
	if *flagShort != defaultValue {
		return flagShort
	}
	
	return flagLong
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "create":
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		
		// تعریف فلگ‌ها
		projectLong := createCmd.String("project", "", "Project name (required)")
		project := createCmd.String("p", "", "Project name (required)")

		nameLong := createCmd.String("name", "", "IME name (required)")
		name := createCmd.String("n", "", "IME name (required)")

		labelLong := createCmd.String("label", "Custom", "IME label (short name)")
		label := createCmd.String("l", "Custom", "IME label (short name)")

		iconLong := createCmd.String("icon", "fcitx-keyboard", "Icon name")
		icon := createCmd.String("i", "fcitx-keyboard", "Icon name")

		langLong := createCmd.String("lang", "en", "Language code")
		lang := createCmd.String("L", "en", "Language code")

		descLong := createCmd.String("desc", "Custom IME", "Description")
		desc := createCmd.String("D", "Custom IME", "Description")

		configLong := createCmd.String("config", "", "Custom config file path")
		config := createCmd.String("c", "", "Custom config file path")

		createCmd.Parse(os.Args[2:])
		
		projectEmpty := (*project == "" && *projectLong == "")
		nameEmpty := (*name == "" && *nameLong == "")
		
		if projectEmpty || nameEmpty {
			fmt.Println("Error: project and name are required")
			createCmd.Usage()
			os.Exit(1)
		}
		
		project = handleDualFlags(project, projectLong, "")
		name = handleDualFlags(name, nameLong, "")
		label = handleDualFlags(label, labelLong, "Custom")
		icon = handleDualFlags(icon, iconLong, "fcitx-keyboard")
		lang = handleDualFlags(lang, langLong, "en")
		desc = handleDualFlags(desc, descLong, "Custom IME")
		config = handleDualFlags(config, configLong, "")

		cfg := Config{
			ProjectName: *project,
			IMEName:     *name,
			Label:       *label,
			Icon:        *icon,
			LangCode:    *lang,
			Description: *desc,
			ConfigFile:  *config,
		}
		
		handleCreate(cfg)

	case "install":
		installCmd := flag.NewFlagSet("install", flag.ExitOnError)
		projectDirLong := installCmd.String("dir", ".", "Project directory")
		projectDir := installCmd.String("d", ".", "Project directory")
		
		installCmd.Parse(os.Args[2:])
		
		projectDir = handleDualFlags(projectDir, projectDirLong, ".")
		cfg := Config{ProjectDir: *projectDir}
		handleInstall(cfg)

	case "edit":
		editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
		nameLong := editCmd.String("name", "", "IME name (required)")
		name := editCmd.String("n", "", "IME name (required)")
		
		editCmd.Parse(os.Args[2:])
		
		name = handleDualFlags(name, nameLong, "")
		
		if *name == "" {
			fmt.Println("Error: name is required")
			editCmd.Usage()
			os.Exit(1)
		}
		
		cfg := Config{IMEName: *name}
		handleEdit(cfg)

	case "list":
		handleList()

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`IME Tool - Custom IME Manager for fcitx5

Usage:
  ime-tool <command> [flags]

Commands:
  create    Create a new IME project
  install   Build and install an IME project
  edit      Edit IME configuration
  list      List installed IMEs

Flags for create:
  -p, --project string    Project name (required)
  -n, --name string       IME name (required)
  -l, --label string      IME label (short name) (default "Custom")
  -i, --icon string       Icon name (default "fcitx-keyboard")
  -L, --lang string       Language code (default "en")
  -D, --desc string       Description (default "Custom IME")
  -c, --config string     Custom config file path

Examples:
  ime-tool create -p myime -n braille -l Brl -c braille.conf
  ime-tool create --project myime --name braille --label Brl --config braille.conf
  ime-tool install -d myime
  ime-tool edit -n braille
  ime-tool list`)
}

func handleCreate(cfg Config) {
	fmt.Printf("Create command called with:\n")
	fmt.Printf("  Project: %s\n", cfg.ProjectName)
	fmt.Printf("  IME Name: %s\n", cfg.IMEName)
	fmt.Printf("  Label: %s\n", cfg.Label)
	fmt.Printf("  Icon: %s\n", cfg.Icon)
	fmt.Printf("  Language: %s\n", cfg.LangCode)
	fmt.Printf("  Config: %s\n", cfg.ConfigFile)
	
	//  createIME(cfg)
}

func handleInstall(cfg Config) {
	fmt.Printf("Install command called for directory: %s\n", cfg.ProjectDir)
	// installIME(cfg) 
}

func handleEdit(cfg Config) {
	fmt.Printf("Edit command called for IME: %s\n", cfg.IMEName)
	// editConfig(cfg) 
}

func handleList() {
	fmt.Println("List command called")
	// listIMEs() 
}