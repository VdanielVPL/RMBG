package image

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"

	_ "golang.org/x/image/webp"

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
	return fileBytes, nil
}

func ToBase64FromBytes(fileBytes []byte) (string, string, error) {
	fileType := http.DetectContentType(fileBytes)
	if strings.HasPrefix(fileType, "image/") {
		return base64.StdEncoding.EncodeToString(fileBytes), fileType, nil
	} else {
		return "", "", errors.New("NOT_AN_IMAGE")
	}
}

func ToBytesFromPath(path string) ([]byte, error) {
	return os.ReadFile(path)
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

func SaveImageDialog(ctx context.Context, isJPG bool) (string, error) {
	options := runtime.SaveDialogOptions{
		Title: "Save an image",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Image (*.png)",
				Pattern:     "*.png",
			},
		},
	}

	if isJPG {
		options.Filters[0].DisplayName = "Image (*.jpg)"
		options.Filters[0].Pattern = "*.jpg"
	}

	var filePath string
	var err error
	done := make(chan struct{})

	go func() {
		filePath, err = runtime.SaveFileDialog(ctx, options)
		close(done)
	}()

	<-done
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func ToJPG(pngImgBytes []byte, jpgImgBytes *[]byte) error {
	img, err := png.Decode(bytes.NewReader(pngImgBytes))
	if err != nil {
		return err
	}

	bounds := img.Bounds()

	rgba := image.NewRGBA(bounds)

	draw.Draw(rgba, bounds, &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	draw.Draw(rgba, bounds, img, bounds.Min, draw.Over)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, rgba, &jpeg.Options{Quality: 90})
	if err != nil {
		return err
	}
	*jpgImgBytes = buf.Bytes()
	return nil
}

func ToPNG(jpgImgBytes []byte, pngImgBytes *[]byte) error {
	img, err := jpeg.Decode(bytes.NewReader(jpgImgBytes))
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return err
	}
	*pngImgBytes = buf.Bytes()
	return nil
}

func GetFormatFromPath(path string) string {
	return path[strings.LastIndex(path, ".")+1:]
}

func GetMimeType(bytes []byte) string {
	return http.DetectContentType(bytes[:512])
}
