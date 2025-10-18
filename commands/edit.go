

func editConfig(imeName string) error {
    // حالا فقط مسیر کاربر رو چک کن
    configPath := filepath.Join(os.Getenv("HOME"), ".config/fcitx5", imeName, "config", imeName+".conf")
    
    // اگه فایل وجود نداره، از مسیر پروژه کپی کن
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        fmt.Printf("Config file not found, creating from project...\n")
        if err := createUserConfigFromProject(imeName); err != nil {
            return err
        }
    }

    fmt.Printf("Editing user config: %s\n", configPath)
    return openEditor(configPath)
}

func createUserConfigFromProject(imeName string) error {
    // پیدا کردن پروژه و کپی کردن کانفیگ به مسیر کاربر
    projectDir := findProjectDir(imeName)
    if projectDir == "" {
        return fmt.Errorf("project directory for IME '%s' not found", imeName)
    }

    configSrc := filepath.Join(projectDir, "config", imeName+".conf")
    configDestDir := filepath.Join(os.Getenv("HOME"), ".config/fcitx5", imeName, "config")
    configDest := filepath.Join(configDestDir, imeName+".conf")
    
    if err := os.MkdirAll(configDestDir, 0755); err != nil {
        return err
    }
    
    return copyFile(configSrc, configDest)
}