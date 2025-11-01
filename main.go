package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"ime-tool/commands"
)

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

func handleDualFlagsBool(flagShort, flagLong *bool) *bool {
	
	if !*flagShort && !*flagLong {
		return flagShort
	}else if *flagLong {
		return  flagLong
	}else if *flagShort {
		return flagShort
	}else {
		
		printUsage()
		os.Exit(1)
		return flagLong
	}
	

	
}

// تابع برای گرفتن تایید از کاربر
func askForConfirmation(question string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [Y/n]: ", question)

		answer, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		// حذف فاصله و کاراکترهای جدید
		answer = strings.TrimSpace(strings.ToLower(answer))

		// اگر Enter زد (پاسخ خالی) یا Y/y زد، تأیید کن
		if answer == "" || answer == "y" || answer == "yes" {
			return true
		}

		// اگر n/no زد، رد کن
		if answer == "n" || answer == "no" {
			return false
		}

		// اگر پاسخ نامعتبر بود، دوباره بپرس
		fmt.Println("❌ Please enter 'y', 'n', or press Enter for yes")
	}
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

		forceLong := createCmd.Bool("force", false, "Force overwrite existing project")
		force := createCmd.Bool("f", false, "Force overwrite existing project")

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
		force = handleDualFlagsBool(force, forceLong)

		// چک کردن کاراکترهای غیرمجاز
		if strings.Contains(*project, "-") || strings.Contains(*name, "-") {
			fmt.Println("Error: Project name and IME name cannot contain '-'")
			os.Exit(1)
		}

		cfg := commands.Config{
			ProjectName: *project,
			IMEName:     *name,
			Label:       *label,
			Icon:        *icon,
			LangCode:    *lang,
			Description: *desc,
			ConfigFile:  *config,
		}

		// اگر force نبود و پروژه وجود داشت، از کاربر سوال بپرس
		if !*force {
			exists, err := commands.ProjectExists(*project)
			if err != nil {
				fmt.Printf("❌ Error checking existing projects: %v\n", err)
				os.Exit(1)
			}

			if exists {
				if !askForConfirmation(fmt.Sprintf("Project '%s' already exists. Overwrite it?", *project)) {
					fmt.Println("Operation canceled.")
					os.Exit(0)
				}
				// اگر کاربر تأیید کرد، force رو true کن
				*force = true
			}
		}

		if err := commands.HandleCreate(cfg, *force); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

	case "install":
		installCmd := flag.NewFlagSet("install", flag.ExitOnError)
		projectNameLong := installCmd.String("project", "", "Project name (required)")
		projectName := installCmd.String("p", "", "Project name (required)")

		installCmd.Parse(os.Args[2:])

		projectName = handleDualFlags(projectName, projectNameLong, "")

		if *projectName == "" {
			fmt.Println("Error: project name is required")
			installCmd.Usage()
			os.Exit(1)
		}

		cfg := commands.Config{ProjectName: *projectName}
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

		cfg := commands.Config{IMEName: *name}
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
  -f, --force             Force overwrite existing project

Examples:
  ime-tool create -p myime -n braille -l Brl -c braille.conf
  ime-tool create --project myime --name braille --label Brl --config braille.conf --force
  ime-tool install -p myime
  ime-tool edit -n braille
  ime-tool list`)
}

func handleInstall(cfg commands.Config) {
	fmt.Printf("Install command called for project: %s\n", cfg.ProjectName)
	// installIME(cfg)
}

func handleEdit(cfg commands.Config) {
	fmt.Printf("Edit command called for IME: %s\n", cfg.IMEName)
	// editConfig(cfg)
}

func handleList() {
	projects, err := commands.HandleList()
	if err != nil {
		fmt.Printf("❌ Error listing projects: %v\n", err)
		return
	}

	if len(projects) == 0 {
		fmt.Println("📋 No IME projects found.")
		fmt.Println("Create your first project with: ime-tool create -p projectname -n imename")
		return
	}

	fmt.Println("📋 Available IME projects:")
	for _, project := range projects {
		fmt.Printf("  • %s\n", project)
	}
}