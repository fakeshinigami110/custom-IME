package main

import (
	"bufio"
	"flag"
	"fmt"
	"ime-tool/commands"
	"os"
	"regexp"
	"strings"
)

func isVaildName(input string) bool {
	pattern := `^[a-zA-z][a-zA-z0-9_]+$`
	is_matched, _ := regexp.MatchString(pattern, input)
	return is_matched
}

func handleDualFlags[T comparable](flagShort, flagLong *T, defaultValue T) *T {
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

func AskForConfirmation(question string) bool {
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

		if projectEmpty == true {
			fmt.Println("Error: project is required")
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
		if isVaildName(*project) == false && isVaildName(*name) == false {
			fmt.Println(
				`Error: Project name and IME name should start with digits and it can contain just latin digits , nmbers and _ `)
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
				if !AskForConfirmation(fmt.Sprintf("Project '%s' already exists. Overwrite it?", *project)) {
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
		projectNameShort := installCmd.String("p", "", "Project name (required)")

		installCmd.Parse(os.Args[2:])

		projectName := handleDualFlags(projectNameShort, projectNameLong, "")

		if *projectName == "" {
			fmt.Println("Error: project name is required")
			installCmd.Usage()
			os.Exit(1)
		}

		cfg := commands.Config{ProjectName: *projectName}
		handleInstall(cfg)

	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		idLong := deleteCmd.Int("id", 0, "enter IME's id that you would like to delete (required)")
		idShort := deleteCmd.Int("i", 0, "enter IME's id that you would like to delete (required)")

		deleteLong := deleteCmd.Bool("d", false, "throw -d or --delete")
		deleteShort := deleteCmd.Bool("delete", false, "throw -d or --delete")

		uninstallShort := deleteCmd.Bool("u", false, "throw -u or --uninstall to uninstall IME from fcitx5")
		uninstallLong := deleteCmd.Bool("uninstall", false, "throw -u or --uninstall to uninstall IME from fcitx5")

		deleteCmd.Parse(os.Args[2:])
		// fmt.Println("hereee2")

		id := handleDualFlags(idLong, idShort, 0)
		if *id == 0 {
			// s , err := commands.GetDBPath()
			// fmt.Printf("filepath : %v , err : %v \n" , s ,err)
			deleteCmd.Usage()
			os.Exit(1)
		}
		// fmt.Println("hereee2")

		del := handleDualFlagsBool(deleteLong, deleteShort)
		un := handleDualFlagsBool(uninstallShort, uninstallLong)
		// fmt.Println("hereee")
		commands.HandleDlete(*id, *del, *un)

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
  delte     use for delete sources or fully unistall an IME.

Flags for create:
  -p, --project string    Project name (required)
  -n, --name string       IME name (required)
  -l, --label string      IME label (short name) (default "Custom")
  -i, --icon string       Icon name (default "fcitx-keyboard")
  -L, --lang string       Language code (default "en")
  -D, --desc string       Description (default "Custom IME")
  -c, --config string     Custom config file path
  -f, --force             Force overwrite existing project

Flags for delete :
  -i, --ID string                IME number (get it from list flag) int    id (required)
  -u, --uninstall             use to uninstall the IME from fcitx5 entries
  -d, --delete                use to delete source codes from your system

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
	projects, _, err := commands.RertuenImes()
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
	for id, project := range projects {
		fmt.Printf("%d %s\n", id, project)
	}
}
