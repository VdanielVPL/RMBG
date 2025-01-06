package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func generateUniqueFileName(dir, baseName string) string {
	suffix := 1
	var newFileName string

	for {
		newFileName = baseName + "-" + strconv.Itoa(suffix) + ".png"
		newPath := filepath.Join(dir, newFileName)

		if !fileExists(newPath) {
			return newPath
		}

		suffix++
	}
}

func Cli(args []string) {
	if args[1] != "" {
		absPath, err := filepath.Abs(args[1])

		if err != nil {
			return
		}

		dir := filepath.Dir(absPath)
		dirName := filepath.Base(dir)
		newFilePath := generateUniqueFileName(dir, dirName)

		if len(args) == 2 {
			cmd := exec.Command("rembg", "i", absPath, newFilePath)
			cmd.SysProcAttr = &syscall.SysProcAttr{
				CreationFlags: 0x08000000,
				HideWindow:    true,
			}
			cmd.Run()
		}

		if len(args) == 3 {
			cmd := exec.Command("rembg", "i", "-m", args[2], absPath, newFilePath)
			cmd.SysProcAttr = &syscall.SysProcAttr{
				CreationFlags: 0x08000000,
				HideWindow:    true,
			}
			cmd.Run()
		}
	}
}
