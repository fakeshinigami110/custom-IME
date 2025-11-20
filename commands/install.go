package commands

import (
	"fmt"
	"os"
	"os/exec"
)

func HandleInstallation (projectName string) error {
	db , err:= LoadDB()
	if err != nil {
		fmt.Printf("%v %v\n" ,ResultColor(false) ,err)
		os.Exit(1)
	}
	 
	cfg , exists := db.GetProject(projectName) 
	if !exists { 
		return  fmt.Errorf("there is no project in name %v" , projectName)
	}
	dir := cfg.ProjectPath
	cmdCompile := exec.Command("cmake" ,"-DCMAKE_INSTALL_PREFIX=/usr")
	cmdCompile.Dir = dir
	cmdCompile.Stderr = os.Stderr
	cmdCompile.Stdin = os.Stdin
	cmdCompile.Stdout = os.Stdout
	err = cmdCompile.Run()
	if err != nil {
		fmt.Printf("%v %v\n" ,ResultColor(false) ,err)
		os.Exit(1)
	}

	cmd := exec.Command("sudo" , "make" , "install")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin =os.Stdin


	err = cmd.Run()
	if err != nil {
		fmt.Printf("%v %v\n" ,ResultColor(false) ,err)
		os.Exit(1)
	}
	cmd2 := exec.Command("fcitx5" , "-rd")
		err = cmd2.Run()
	if err != nil {
		fmt.Printf("%v %v\n" ,ResultColor(false) ,err)
		os.Exit(1)
	}
	fmt.Println("isntallation complated")
	cfg.Installed = true
	db.UpdateProject(*cfg)
	return  err

	

	
}