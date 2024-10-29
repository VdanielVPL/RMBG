package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) > 1 {
		absPath := os.Args[1]
		dir := filepath.Dir(absPath)
		dirName := filepath.Base(dir)
		newFilePath := generateUniqueFileName(dir, dirName)
		if len(os.Args) == 2 {
			cmd := exec.Command("rembg", "i", absPath, newFilePath)
			cmd.SysProcAttr = &syscall.SysProcAttr{
				HideWindow: true,
			}
			cmd.Run()
		}
		if len(os.Args) == 3 {
			cmd := exec.Command("rembg", "i", "-m", os.Args[2], absPath, newFilePath)
			cmd.SysProcAttr = &syscall.SysProcAttr{
				HideWindow: true,
			}
			cmd.Run()
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func generateUniqueFileName(dir, baseName string) string {
	suffix := 1
	var newFileName string
	for {
		// Tworzymy nową nazwę z sufiksem
		newFileName = baseName + "-" + strconv.Itoa(suffix) + ".png"
		newPath := filepath.Join(dir, newFileName)

		// Sprawdzamy, czy plik o tej nazwie już istnieje
		if !fileExists(newPath) {
			return newPath
		}
		suffix++
	}
}
