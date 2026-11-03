package utils

import (
	"context"
	"os"
	"path/filepath"
	"rmbg/utils/image"
	"strconv"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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

func Cli(args []string) error {
	if args[1] == "" {
		return nil
	}

	absPath, err := filepath.Abs(args[1])
	if err != nil {
		return err
	}

	dir := filepath.Dir(absPath)
	dirName := filepath.Base(dir)
	newFilePath := generateUniqueFileName(dir, dirName)

	rembgArgs := []string{"i", absPath, newFilePath}
	if len(args) > 2 {
		rembgArgs = []string{"i", "-m", args[2], absPath, newFilePath}
	}

	cmd := image.RembgCommand(context.Background(), rembgArgs...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
