package image

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/disintegration/imaging"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var supportedImageTypes = []string{"png", "jpeg", "webp"}

func IsSupportedImage(fileBytes []byte) bool {
	_, format, err := image.DecodeConfig(bytes.NewReader(fileBytes))
	return err == nil && slices.Contains(supportedImageTypes, format)
}

const maxURLImageSize = 100 << 20 // 100 MB

var httpClient = &http.Client{Timeout: 60 * time.Second}

func GetImageFromURL(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "RMBG") // some servers (e.g. Wikipedia) reject Go's default one

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed: %s", resp.Status)
	}

	fileBytes, err := io.ReadAll(io.LimitReader(resp.Body, maxURLImageSize+1))
	if err != nil {
		return nil, err
	}
	if len(fileBytes) > maxURLImageSize {
		return nil, errors.New("downloaded image is over 100 MB")
	}
	return fileBytes, nil
}

func ToBase64FromBytes(fileBytes []byte) (string, string, error) {
	fileType := http.DetectContentType(fileBytes)
	if strings.HasPrefix(fileType, "image/") {
		return base64.StdEncoding.EncodeToString(fileBytes), fileType, nil
	} else {
		return "", "", errors.New("NO_IMAGE")
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
	return runtime.OpenFileDialog(ctx, options)
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

	return runtime.SaveFileDialog(ctx, options)
}

func ToJPG(pngImgBytes []byte, jpgImgBytes *[]byte) error {
	img, _, err := decodeOriented(pngImgBytes)
	if err != nil {
		return err
	}

	bounds := img.Bounds()

	rgba := image.NewRGBA(bounds)

	draw.Draw(rgba, bounds, &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	draw.Draw(rgba, bounds, img, bounds.Min, draw.Over)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, rgba, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}
	*jpgImgBytes = buf.Bytes()
	return nil
}

func ToPNG(jpgImgBytes []byte, pngImgBytes *[]byte) error {
	img, _, err := decodeOriented(jpgImgBytes)
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

// JPG and PNG stay untouched, other formats are converted to lossless PNG
func ToSupportedFormat(imgBytes []byte) (newImgBytes []byte, isJPG bool, converted bool, err error) {
	switch GetMimeType(imgBytes) {
	case "image/jpeg":
		return imgBytes, true, false, nil
	case "image/png":
		return imgBytes, false, false, nil
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, false, false, errors.New("INVALID_IMAGE")
	}
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, false, false, err
	}
	return buf.Bytes(), false, true, nil
}

func decodeOriented(imgBytes []byte) (image.Image, string, error) {
	img, format, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, format, err
	}

	if format == "jpeg" {
		img = ApplyOrientation(img, ReadOrientation(bytes.NewReader(imgBytes)))
	}
	return img, format, nil
}

// ApplyOrientation turns the image the way EXIF orientation says, like the preview shows it
func ApplyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.FlipV(img)
	case 5:
		return imaging.Transpose(img)
	case 6:
		return imaging.Rotate270(img)
	case 7:
		return imaging.Transverse(img)
	case 8:
		return imaging.Rotate90(img)
	}
	return img
}

func GetMimeType(bytes []byte) string {
	return http.DetectContentType(bytes)
}
