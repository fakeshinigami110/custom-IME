package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

func HandleErr (err error){
	if err != nil {
		fmt.Printf("something went wrong : %v" , err)
		os.Exit(1)
	}
}

func clearDirctory (dirpath string)  {
	if _ ,err := os.Stat(dirpath) ; os.IsNotExist(err){
		HandleErr(err)
	}
	entries , err := os.ReadDir(dirpath)
	HandleErr(err)
	for _ ,entry := range entries {
		fp := filepath.Join(dirpath , entry.Name())

		if entry.IsDir() {
			err := os.RemoveAll(fp)
			HandleErr(err)

		}else {
			err := os.Remove(fp)
			HandleErr(err)
		}
	} 
}


func HandleDlete(id int , del , un bool){
	imes ,baseDir, err:= RertuenImes()
	if err != nil {
		fmt.Printf("An error accourd : %v" , err)
		os.Exit(1)
	}
	value , exists := imes[id]
	dir := filepath.Join(baseDir,value)
	if exists == false {
		fmt.Printf("there is no ime with id %d" , id)

	}
	if del == true {
		clearDirctory(dir)
	}
	if un == true {
		// I LEFT HERE 
	}
	
}

// delete source codes in 
/*
~/.config/custom-ime/*
*/


// uninstall IME: we should remove the following filese and restart via fcix5 -rd
/*
Install the project...
-- Install configuration: "Debug"
-- Installing: /usr/share/fcitx5/cuneiform_ime/config/cuneiform_ime.conf
-- Up-to-date: /usr/lib/fcitx5/cuneiform_ime.so
-- Up-to-date: /usr/share/fcitx5/addon/cuneiform_ime.conf
-- Up-to-date: /usr/share/fcitx5/inputmethod/cuneiform_ime.conf

*/