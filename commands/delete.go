package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
)

func HandleErr(err error) {
	if err != nil {
		fmt.Printf("%v something went wrong : %v\n",ResultColor(false), err)
		// os.Exit(1)
	}
}
func Check_dir(dirpath string) {
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		HandleErr(err)
	}
}

func clearDirctory(dirpath string) {
	Check_dir(dirpath)
	entries, err := os.ReadDir(dirpath)
	HandleErr(err)
	for _, entry := range entries {
		fp := filepath.Join(dirpath, entry.Name())

		if entry.IsDir() {
			clearDirctory(fp)

		} else {
			err := os.Remove(fp)
			HandleErr(err)
		}
		fmt.Printf("removing %v\n", fp)
	}
}

func HandleDelete(projectName string, del, un bool) {
	db , err := LoadDB()
	cfg , _ := db.GetProject(projectName)
	imes, baseDir, err := RertuenImes()
	projects := db.ListProjects()
	if err != nil {
		fmt.Printf("%v An error accourd : %v", ResultColor(false),err)
		os.Exit(1)
	}
	if !(slices.Contains(imes , projectName) || slices.Contains(projects , projectName)){
		// fmt.Println(projects , projectName)
		// fmt.Println(slices.Contains(projects , projectName))
		fmt.Printf("%v there is no IME in name %v\n" , ResultColor(false),projectName)
		os.Exit(1)
	}
	
	dir := filepath.Join(baseDir, projectName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if del == true {
		clearDirctory(dir)
		os.RemoveAll(dir)
		fmt.Printf("%v removed : %v\n" , ResultColor(true),dir)
		if cfg , exists := db.GetProject(projectName) ; exists {
			cfg.Removed = true
			err := db.UpdateProject(*cfg)
			err = db.Save()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("%v there is no IME in name %v\n" , ResultColor(false),projectName)
			// os.Exit(1)
		}
		cfg.Removed = true
		fmt.Printf("%v source codes of IME %v successfully removed\n" , ResultColor(true) , projectName)



	}
if un == true {
    ImesBaseDir := filepath.Join("/usr", "share", "fcitx5")
    
    targetDir := filepath.Join(ImesBaseDir, projectName)
    cmd := exec.Command("sudo", "rm", "-rf", targetDir)
    err := cmd.Run()
    HandleErr(err)
    fmt.Printf("removing %v\n", targetDir)

    file := filepath.Join("/usr", "lib", "fcitx5", projectName+".so")
    cmd = exec.Command("sudo", "rm", "-f", file)
    err = cmd.Run()
    HandleErr(err)
    fmt.Printf("removing %v\n", file)

    file = filepath.Join(ImesBaseDir, "addon", fmt.Sprintf("%v.conf", projectName))
    cmd = exec.Command("sudo", "rm", "-f", file)
    err = cmd.Run()
    HandleErr(err)
    fmt.Printf("removing %v\n", file)

    file = filepath.Join(ImesBaseDir, "inputmethod", fmt.Sprintf("%v.conf", projectName))
    cmd = exec.Command("sudo", "rm", "-f", file)
    err = cmd.Run()
    HandleErr(err)
    fmt.Printf("removing %v\n", file)

    fmt.Println("updating fcitx5 . . .")	
    cmd = exec.Command("fcitx5", "-rd")
    err = cmd.Run()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
	cfg.Installed = false
	fmt.Printf("%v IME %v successfully uninstalled\n" , ResultColor(true) , projectName)
	
	


db.UpdateProject(*cfg)
db.Save()
}
    
if !cfg.Installed && cfg.Removed {
	err = db.DeleteProject(projectName)
	if err != nil {
	fmt.Println(err)
	}
}
}

// delete source codes in
/*
~/.config/custom-ime/*
*/

// uninstall IME: we should remove the following filese and restart via fcix5 -rd
/*
Install the project...
-- Install configuration: ""
-- Installing: /usr/share/fcitx5/cuneiform_imeP/config/cuneiform_imeN.conf
-- Installing: /usr/lib/fcitx5/cuneiform_imeN.so
-- Installing: /usr/share/fcitx5/addon/cuneiform_imeN.conf
-- Installing: /usr/share/fcitx5/inputmethod/cuneiform_imeN.conf

*/
