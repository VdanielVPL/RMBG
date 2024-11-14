package image

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"
	"unsafe"

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

// What is this...
func CopyToClipboard(imageData []byte) error {
	// Dekoduj obraz z bajtów
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return err
	}
	fmt.Println("Format obrazu:", format)

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Przygotuj dane pikseli w formacie BGRA z kanałem alfa
	pixelData := make([]byte, width*height*4)
	idx := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Pobierz kolor piksela
			r, g, b, a := img.At(x, y).RGBA()
			// Konwertuj wartości z zakresu 0-65535 do 0-255
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// Zapisz piksel w formacie BGRA
			pixelData[idx+0] = b8
			pixelData[idx+1] = g8
			pixelData[idx+2] = r8
			pixelData[idx+3] = a8
			idx += 4
		}
	}

	// Tworzenie nagłówka BITMAPV5HEADER
	var biHeader BITMAPV5HEADER
	biHeader.BV5Size = uint32(unsafe.Sizeof(biHeader))
	biHeader.BV5Width = int32(width)
	biHeader.BV5Height = int32(-height) // Ujemna wartość dla obrazu top-down
	biHeader.BV5Planes = 1
	biHeader.BV5BitCount = 32
	biHeader.BV5Compression = BI_BITFIELDS
	biHeader.BV5RedMask = 0x00FF0000
	biHeader.BV5GreenMask = 0x0000FF00
	biHeader.BV5BlueMask = 0x000000FF
	biHeader.BV5AlphaMask = 0xFF000000
	biHeader.BV5CSType = LCS_sRGB

	// Połącz nagłówek z danymi pikseli
	var dib bytes.Buffer
	binary.Write(&dib, binary.LittleEndian, &biHeader)
	dib.Write(pixelData)

	// Zakoduj obraz do PNG
	var pngBuffer bytes.Buffer
	err = png.Encode(&pngBuffer, img)
	if err != nil {
		return err
	}
	pngData := pngBuffer.Bytes()

	// Ustawienie obrazu w schowku
	err = setClipboardImage(dib.Bytes(), pngData)
	if err != nil {
		return err
	}

	fmt.Println("Obraz został skopiowany do schowka z przezroczystością.")
	return nil
}

// Definicje struktur i stałych

const (
	CF_DIBV5      = 17
	BI_BITFIELDS  = 3
	LCS_sRGB      = 0x73524742
	GMEM_MOVEABLE = 0x0002
)

var CF_PNG uint32

func init() {
	formatName, _ := syscall.UTF16PtrFromString("PNG")
	CF_PNG = RegisterClipboardFormat(formatName)
}

func RegisterClipboardFormat(formatName *uint16) uint32 {
	user32 := syscall.NewLazyDLL("user32.dll")
	registerClipboardFormat := user32.NewProc("RegisterClipboardFormatW")
	ret, _, _ := registerClipboardFormat.Call(uintptr(unsafe.Pointer(formatName)))
	return uint32(ret)
}

type BITMAPV5HEADER struct {
	BV5Size          uint32
	BV5Width         int32
	BV5Height        int32
	BV5Planes        uint16
	BV5BitCount      uint16
	BV5Compression   uint32
	BV5SizeImage     uint32
	BV5XPelsPerMeter int32
	BV5YPelsPerMeter int32
	BV5ClrUsed       uint32
	BV5ClrImportant  uint32
	BV5RedMask       uint32
	BV5GreenMask     uint32
	BV5BlueMask      uint32
	BV5AlphaMask     uint32
	BV5CSType        uint32
	BV5Endpoints     CIEXYZTRIPLE
	BV5GammaRed      uint32
	BV5GammaGreen    uint32
	BV5GammaBlue     uint32
	BV5Intent        uint32
	BV5ProfileData   uint32
	BV5ProfileSize   uint32
	BV5Reserved      uint32
}

type CIEXYZTRIPLE struct {
	CiexyzRed   CIEXYZ
	CiexyzGreen CIEXYZ
	CiexyzBlue  CIEXYZ
}

type CIEXYZ struct {
	CiexyzX int32
	CiexyzY int32
	CiexyzZ int32
}

func setClipboardImage(dibData []byte, pngData []byte) error {
	user32 := syscall.NewLazyDLL("user32.dll")
	kernel32 := syscall.NewLazyDLL("kernel32.dll")

	openClipboard := user32.NewProc("OpenClipboard")
	emptyClipboard := user32.NewProc("EmptyClipboard")
	setClipboardData := user32.NewProc("SetClipboardData")
	closeClipboard := user32.NewProc("CloseClipboard")

	globalAlloc := kernel32.NewProc("GlobalAlloc")
	globalLock := kernel32.NewProc("GlobalLock")
	globalUnlock := kernel32.NewProc("GlobalUnlock")

	// Otwórz schowek
	ret, _, err := openClipboard.Call(0)
	if ret == 0 {
		return err
	}
	defer closeClipboard.Call()

	// Wyczyść schowek
	ret, _, err = emptyClipboard.Call()
	if ret == 0 {
		return err
	}

	// Ustaw dane w formacie CF_DIBV5
	err = setClipboardDataFormat(globalAlloc, globalLock, globalUnlock, setClipboardData, CF_DIBV5, dibData)
	if err != nil {
		fmt.Println("Błąd ustawiania CF_DIBV5:", err)
	}

	// Ustaw dane w formacie CF_PNG
	err = setClipboardDataFormat(globalAlloc, globalLock, globalUnlock, setClipboardData, CF_PNG, pngData)
	if err != nil {
		fmt.Println("Błąd ustawiania CF_PNG:", err)
	}

	return nil
}

func setClipboardDataFormat(globalAlloc, globalLock, globalUnlock, setClipboardData *syscall.LazyProc, format uint32, data []byte) error {
	hMem, _, err := globalAlloc.Call(GMEM_MOVEABLE, uintptr(len(data)))
	if hMem == 0 {
		return err
	}

	memPtr, _, err := globalLock.Call(hMem)
	if memPtr == 0 {
		return err
	}
	defer globalUnlock.Call(hMem)

	copy((*[1 << 30]byte)(unsafe.Pointer(memPtr))[:len(data):len(data)], data)

	ret, _, err := setClipboardData.Call(uintptr(format), hMem)
	if ret == 0 {
		return err
	}

	return nil
}
