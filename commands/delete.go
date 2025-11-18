package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

func HandleErr(err error) {
	if err != nil {
		fmt.Printf("something went wrong : %v", err)
		os.Exit(1)
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

func HandleDlete(id int, del, un bool) {
	imes, baseDir, err := RertuenImes()
	if err != nil {
		fmt.Printf("An error accourd : %v", err)
		os.Exit(1)
	}
	value, exists := imes[id]
	dir := filepath.Join(baseDir, value)
	if exists == false {
		fmt.Printf("there is no ime with id %d", id)

	}
	if del == true {
		clearDirctory(dir)
	}
	if un == true {
		ImesBaseDir := filepath.Join("usr", "share", "fcitx5")
		Check_dir(ImesBaseDir)
		// removes exp : /usr/share/fcitx5/cuneiform_ime/config/cuneiform_ime.conf
		targetDir := filepath.Join(ImesBaseDir, imes[id])
		err := os.RemoveAll(targetDir)
		HandleErr(err)
		fmt.Printf("removing %v", targetDir)

		// removes : /usr/lib/fcitx5/cuneiform_ime.so
		file := filepath.Join("usr", "lib", "fcitx5", imes[id]+".so")

		Check_dir(file)

		err = os.Remove(file)
		HandleErr(err)
		fmt.Printf("removing %v", file)

		// /usr/share/fcitx5/addon/cuneiform_ime.conf
		file = filepath.Join(ImesBaseDir, "addon", fmt.Sprintf("%v.conf", imes[id]))

		Check_dir(file)
		err = os.Remove(file)
		HandleErr(err)
		fmt.Printf("removing %v", file)

		// Up-to-date: /usr/share/fcitx5/inputmethod/cuneiform_ime.conf

		file = filepath.Join(ImesBaseDir, "inputmethod", fmt.Sprintf("%v.conf", imes[id]))

		Check_dir(file)
		err = os.Remove(file)
		HandleErr(err)
		fmt.Printf("removing %v", file)

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
