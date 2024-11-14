package image

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func ToBase64FromPath(path string) (string, string, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
	}
	return ToBase64FromBytes(fileBytes)
}
func GetImageFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return fileBytes, err
}

func ToBase64FromBytes(fileBytes []byte) (string, string, error) {
	fileType := http.DetectContentType(fileBytes)
	if strings.HasPrefix(fileType, "image/") {
		return base64.StdEncoding.EncodeToString(fileBytes), fileType, nil
	} else {
		return "", "", nil
	}
}

func OpenImageDialog(ctx context.Context) (string, error) {
	options := runtime.OpenDialogOptions{
		Title: "Select an image",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Images (*.png, *.jpg, *.jpeg, *.webp)",
				Pattern:     "*.png;*.jpg;*.jpeg;*.webp",
			},
		},
	}
	var filePath string
	var err error
	done := make(chan struct{})

	go func() {
		filePath, err = runtime.OpenFileDialog(ctx, options)
		close(done)
	}()

	<-done
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func SaveImageDialog(ctx context.Context) (string, error) {
	options := runtime.SaveDialogOptions{
		Title: "Save an image",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Image (*.png)",
				Pattern:     "*.png",
			},
		},
	}
	var filePath string
	var errr error
	done := make(chan struct{})

	go func() {
		filePath, errr = runtime.SaveFileDialog(ctx, options)
		close(done)
	}()

	<-done
	if errr != nil {
		return "", errr
	}
	return filePath, nil
}
