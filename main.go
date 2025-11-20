package main

import (
	"bufio"
	"flag"
	"fmt"
	"ime-tool/commands"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func isValidName(input string) bool {
	pattern := `^[a-zA-Z][a-zA-Z0-9_]+$`
	isMatched, _ := regexp.MatchString(pattern, input)
	return isMatched
}

func handleDualFlags[T comparable](flagShort, flagLong *T, defaultValue T) *T {
	if *flagShort == defaultValue && *flagLong == defaultValue {
		return flagShort
	}

	if *flagShort != defaultValue && *flagLong != defaultValue && *flagShort != *flagLong {
		fmt.Printf("%v ", commands.ResultColor(false))
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
	} else if *flagLong {
		return flagLong
	} else if *flagShort {
		return flagShort
	} else {
		printUsage()
		os.Exit(1)
		return flagLong
	}
}

func askForConfirmation(question string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [Y/n]: ", question)

		answer, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		answer = strings.TrimSpace(strings.ToLower(answer))

		if answer == "" || answer == "y" || answer == "yes" {
			return true
		}

		if answer == "n" || answer == "no" {
			return false
		}
		fmt.Printf("%v ", commands.ResultColor(false))
		fmt.Println("Please enter 'y', 'n', or press Enter for yes")
	}
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "create":
		handleCreateCommand()
	case "install":
		handleInstallCommand()
	case "delete":
		handleDeleteCommand()
	case "list":
		handleList()
	default:
		printUsage()
	}
}

func handleCreateCommand() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)

	// Project flags
	projectLong := createCmd.String("project", "", "Project name (required)")
	project := createCmd.String("p", "", "Project name (required)")

	// IME flags
	nameLong := createCmd.String("name", "", "IME name (required)")
	name := createCmd.String("n", "", "IME name (required)")

	// Optional flags with defaults
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

	// Validate required flags
	projectEmpty := (*project == "" && *projectLong == "")
	if projectEmpty {
		fmt.Printf("%v Error: project name is required\n", commands.ResultColor(false))
		createCmd.Usage()
		os.Exit(1)
	}

	// Handle dual flags
	project = handleDualFlags(project, projectLong, "")
	name = handleDualFlags(name, nameLong, "")
	label = handleDualFlags(label, labelLong, "Custom")
	icon = handleDualFlags(icon, iconLong, "fcitx-keyboard")
	lang = handleDualFlags(lang, langLong, "en")
	desc = handleDualFlags(desc, descLong, "Custom IME")
	config = handleDualFlags(config, configLong, "")
	force = handleDualFlagsBool(force, forceLong)

	// Validate names
	if !isValidName(*project) || !isValidName(*name) {
		fmt.Printf("%v ", commands.ResultColor(false))
		fmt.Println("Error: Project name and IME name should start with a letter and can only contain letters, numbers, and underscores")
		os.Exit(1)
	}

	// Handle existing project confirmation
	if !*force {
		exists, err := commands.ProjectExists(*project)
		if err != nil {
			fmt.Printf("%v Error checking existing projects: %v\n", commands.ResultColor(false), err)
			os.Exit(1)
		}

		if exists {
			if !askForConfirmation(fmt.Sprintf("Project '%s' already exists. Overwrite it?", *project)) {
				fmt.Println("Operation canceled.")
				os.Exit(0)
			}
			*force = true
		}
	}

	// Create configuration and execute
	cfg := commands.Config{
		ProjectName: *project,
		IMEName:     *name,
		Label:       *label,
		Icon:        *icon,
		LangCode:    *lang,
		Description: *desc,
		ConfigFile:  *config,
	}

	if err := commands.HandleCreate(cfg, *force); err != nil {
		fmt.Printf("%v Error: %v\n", commands.ResultColor(false), err)
		os.Exit(1)
	}
}

func handleInstallCommand() {
	installCmd := flag.NewFlagSet("install", flag.ExitOnError)
	projectNameLong := installCmd.String("project", "", "Project name (required)")
	projectNameShort := installCmd.String("p", "", "Project name (required)")

	installCmd.Parse(os.Args[2:])

	projectName := handleDualFlags(projectNameShort, projectNameLong, "")

	if *projectName == "" {
		fmt.Printf("%v Error: project name is required\n", commands.ResultColor(false))
		installCmd.Usage()
		os.Exit(1)
	}

	err := commands.HandleInstallation(*projectName)
	if err != nil {
		fmt.Printf("%v %v\n", commands.ResultColor(false), err)
		os.Exit(1)
	}
}

func handleDeleteCommand() {
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	
	// Project identification
	idLong := deleteCmd.String("project", "", "Project name to delete (required)")
	idShort := deleteCmd.String("p", "", "Project name to delete (required)")

	// Action flags
	deleteLong := deleteCmd.Bool("d", false, "Delete source files from system")
	deleteShort := deleteCmd.Bool("delete", false, "Delete source files from system")

	uninstallShort := deleteCmd.Bool("u", false, "Uninstall IME from fcitx5")
	uninstallLong := deleteCmd.Bool("uninstall", false, "Uninstall IME from fcitx5")

	deleteCmd.Parse(os.Args[2:])

	id := handleDualFlags(idLong, idShort, "")
	if *id == "" {
		fmt.Printf("%v Error: project name is required\n", commands.ResultColor(false))
		deleteCmd.Usage()
		os.Exit(1)
	}

	del := handleDualFlagsBool(deleteLong, deleteShort)
	un := handleDualFlagsBool(uninstallShort, uninstallLong)
	
	if !*del && !*un {
		fmt.Printf("%v Error: you must specify at least one action (-d/--delete or -u/--uninstall)\n", commands.ResultColor(false))
		deleteCmd.Usage()
		os.Exit(1)
	}

	commands.HandleDelete(*id, *del, *un)
}

func handleList() {
	db, err := commands.LoadDB()
	if err != nil {
		fmt.Printf("%v Error loading database: %v\n", commands.ResultColor(false), err)
		os.Exit(1)
	}

	projects := db.ListProjects()
	if len(projects) == 0 {
		fmt.Println("No IME projects found.")
		fmt.Println("Create your first project with: custom-ime create -p projectname -n imename")
		return
	}
	

	fmt.Println("Available IME projects:")
	for id, project := range projects {
		cfg , _ := db.GetProject(project)
		fmt.Printf("%s. %s | is Installed : %v | source codes has been removed : %v\n", commands.ColorIt(strconv.Itoa(id+1), "blue"), project , cfg.Installed , cfg.Removed)
	}
}

func printUsage() {
	fmt.Println(`custom-ime - Custom IME Manager for fcitx5

USAGE:
  custom-ime <command> [flags]

COMMANDS:
  create    Create a new IME project
  install   Build and install an IME project
  list      List all IME projects
  delete    Delete sources or uninstall an IME

CREATE COMMAND:
  custom-ime create [flags]

  Required Flags:
    -p, --project string    Project name (letters, numbers, underscore only)
    -n, --name string       IME name (letters, numbers, underscore only)

  Optional Flags:
    -l, --label string      IME label (short name) (default: "Custom")
    -i, --icon string       Icon name (default: "fcitx-keyboard")
    -L, --lang string       Language code (default: "en")
    -D, --desc string       Description (default: "Custom IME")
    -c, --config string     Custom config file path
    -f, --force             Force overwrite existing project

INSTALL COMMAND:
  custom-ime install [flags]

  Required Flags:
    -p, --project string    Project name to install

DELETE COMMAND:
  custom-ime delete [flags]

  Required Flags:
    -p, --project string    Project name to operate on

  Action Flags (at least one required):
    -d, --delete            Delete source files from system
    -u, --uninstall         Uninstall IME from fcitx5

EXAMPLES:
  # Create a new IME project
  custom-ime create -p myime -n braille -l Brl -c braille.conf
  custom-ime create --project myime --name braille --label Brl --config braille.conf --force

  # Install an IME project
  custom-ime install -p myime
  custom-ime install --project myime

  # Delete operations
  custom-ime delete -p myime -d              # Delete source files only
  custom-ime delete -p myime -u              # Uninstall from fcitx5 only  
  custom-ime delete -p myime -d -u           # Delete sources and uninstall

  # List all projects
  custom-ime list`)
}