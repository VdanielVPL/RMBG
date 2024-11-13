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

// What is this...
func CopyToClipboard(rembgimage []byte) error {
	// if rembgimage != nil {
	// 	// Zdekoduj obraz z bajtów
	// 	img, format, err := image.Decode(bytes.NewReader(rembgimage))
	// 	if err != nil {
	// 		fmt.Println("Błąd dekodowania obrazu:", err)
	// 		return err
	// 	}
	// 	fmt.Println("Format obrazu:", format)

	// 	// Utwórz nowy obraz RGBA o tych samych rozmiarach
	// 	bounds := img.Bounds()
	// 	rgba := image.NewRGBA(bounds)

	// 	// Przejdź przez każdy piksel
	// 	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	// 		for x := bounds.Min.X; x < bounds.Max.X; x++ {
	// 			// Pobierz kolor piksela
	// 			origColor := img.At(x, y)
	// 			r, g, b, a := origColor.RGBA()

	// 			// Konwertuj wartości z zakresu 0-65535 do 0-255
	// 			r8 := uint8(r >> 8)
	// 			g8 := uint8(g >> 8)
	// 			b8 := uint8(b >> 8)
	// 			a8 := uint8(a >> 8)

	// 			if a8 == 0 {
	// 				// Jeśli piksel jest całkowicie przezroczysty, ustaw na biały
	// 				rgba.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	// 			} else {
	// 				// W przeciwnym razie zachowaj oryginalny kolor z pełną nieprzezroczystością
	// 				rgba.Set(x, y, color.RGBA{R: r8, G: g8, B: b8, A: 255})
	// 			}
	// 		}
	// 	}

	// 	// Zakoduj zmodyfikowany obraz do formatu PNG
	// 	var buf bytes.Buffer
	// 	if err := png.Encode(&buf, rgba); err != nil {
	// 		fmt.Println("Błąd kodowania obrazu:", err)
	// 		return err
	// 	}

	// 	// Skopiuj obraz do schowka lub zapisz do pliku
	// 	// Przykład kopiowania do schowka:
	// 	if err := clipboard.Write(clipboard.FmtImage, buf.Bytes()); err != nil {
	// 		fmt.Println("Błąd kopiowania do schowka:", err)
	// 		return nil
	// 	}

	// 	fmt.Println("Obraz został skopiowany do schowka bez przezroczystości.")
	// 	return nil
	// }
	// return nil
	return nil
}
