package commands

import (
	"os"
	"path/filepath"
)

func RertuenImes() (map[int]string ,string , error){
	baseDir := filepath.Join(os.Getenv("HOME"), ".config", "custom-ime")
	counter := 1
	projects := make(map[int]string)

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return projects, baseDir,nil 
		}
		return nil, baseDir,err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			projects[counter] = entry.Name()
			counter++
		}
	}
	return projects , baseDir,err

}

// func HandleList() ([]string, error) {
// 	baseDir := filepath.Join(os.Getenv("HOME"), ".config", "custom-ime")

// 	var projects []string

// 	entries, err := os.ReadDir(baseDir)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return projects, nil // دایرکتوری وجود نداره = هیچ پروژه‌ای نیست
// 		}
// 		return nil, err
// 	}

// 	for _, entry := range entries {
// 		if entry.IsDir() {
// 			projects = append(projects, entry.Name())
// 		}
// 	}

// 	return projects, nil
// }

func ProjectExists(projectName string) (bool, error) {
	projects, _,err := RertuenImes()
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