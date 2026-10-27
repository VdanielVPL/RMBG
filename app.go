package main

import (
	"context"
	"os"
	"rmbg/utils"
	"rmbg/utils/image"
	"rmbg/utils/image/clipboard"
	goruntime "runtime"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	LangStrings map[string]string
	IsDarkMode  bool
	AccentColor string
	imgData     *imgDataStruct

	rmbgMu sync.Mutex // not held while rembg runs
	cropMu sync.Mutex // never held together with rmbgMu
}

type imgDataStruct struct {
	rembgpath      string
	rembgimg       []byte
	rembgimg2      []byte
	rembgGen       int // changes with every new RMBG input
	model          string
	pngcropimgpath string
	pngcropimg     []byte
	jpgcropimgpath string
	jpgcropimg     []byte
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if langData, err := utils.LoadLang(); err == nil {
		a.LangStrings = langData
	} else {
		runtime.EventsEmit(ctx, "alert", err.Error())
	}

	if accentColor, err := utils.GetAccentColor(); err == nil {
		a.AccentColor = accentColor
	} else {
		runtime.EventsEmit(ctx, "alert", err.Error())
	}

	a.IsDarkMode = utils.IsDarkMode()

	a.imgData = &imgDataStruct{
		model: "u2net",
	}

	if err := clipboard.Init(); err != nil {
		runtime.EventsEmit(ctx, "alert", err.Error())
	}
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	clipboard.Close()
}

func (a *App) GetLangStrings() map[string]string {
	return a.LangStrings
}

func (a *App) GetDarkMode() bool {
	return a.IsDarkMode
}

func (a *App) GetAccentColor() string {
	return a.AccentColor
}

func (a *App) HandleDrop(Imagetype string, path string, isUrl bool) []string {
	var imgBytes []byte
	var err error

	if isUrl {
		imgBytes, err = image.GetImageFromURL(path)
		path = ""
	} else {
		imgBytes, err = image.ToBytesFromPath(path)
	}
	if err != nil {
		runtime.EventsEmit(a.ctx, "alert", err.Error())
		return nil
	}

	return a.setInputImage(Imagetype, path, imgBytes)
}

func (a *App) OpenImage(Imagetype string) []string {
	path, err := image.OpenImageDialog(a.ctx)
	if err != nil || path == "" {
		return nil
	}

	imgBytes, err := image.ToBytesFromPath(path)
	if err != nil {
		runtime.EventsEmit(a.ctx, "alert", err.Error())
		return nil
	}

	return a.setInputImage(Imagetype, path, imgBytes)
}

func (a *App) setInputImage(Imagetype string, path string, imgBytes []byte) []string {
	if !image.IsSupportedImage(imgBytes) {
		runtime.EventsEmit(a.ctx, "alert", "NO_IMAGE")
		return nil
	}

	if Imagetype == "RMBG" {
		a.rmbgMu.Lock()
		defer a.rmbgMu.Unlock()
		a.setRMBGImage(path, imgBytes)
	} else if Imagetype == "CROP" {
		a.cropMu.Lock()
		defer a.cropMu.Unlock()
		var err error
		imgBytes, err = a.setCropImage(path, imgBytes)
		if err != nil {
			runtime.EventsEmit(a.ctx, "alert", err.Error())
			return nil
		}
	}

	str, fileType, err := image.ToBase64FromBytes(imgBytes)
	if err != nil {
		runtime.EventsEmit(a.ctx, "alert", err.Error())
		return nil
	}
	return []string{str, fileType}
}

func (a *App) setRMBGImage(path string, imgBytes []byte) {
	a.imgData.rembgpath = path
	a.imgData.rembgimg = nil
	if path == "" {
		a.imgData.rembgimg = imgBytes
	}
	a.imgData.rembgGen++
}

// JPG and PNG files are kept as paths, other formats are converted to PNG and kept in memory
func (a *App) setCropImage(path string, imgBytes []byte) ([]byte, error) {
	imgBytes, isJPG, converted, err := image.ToSupportedFormat(imgBytes)
	if err != nil {
		return nil, err
	}
	if converted {
		path = ""
	}

	var inMemory []byte
	if path == "" {
		inMemory = imgBytes
	}

	if isJPG {
		a.imgData.jpgcropimgpath, a.imgData.jpgcropimg = path, inMemory
		a.imgData.pngcropimgpath, a.imgData.pngcropimg = "", nil
	} else {
		a.imgData.pngcropimgpath, a.imgData.pngcropimg = path, inMemory
		a.imgData.jpgcropimgpath, a.imgData.jpgcropimg = "", nil
	}
	return imgBytes, nil
}

func (a *App) cropImageBytes(isJPG bool) ([]byte, error) {
	if isJPG {
		return readImage(a.imgData.jpgcropimg, a.imgData.jpgcropimgpath)
	}
	return readImage(a.imgData.pngcropimg, a.imgData.pngcropimgpath)
}

func readImage(imgBytes []byte, path string) ([]byte, error) {
	if imgBytes != nil || path == "" {
		return imgBytes, nil
	}
	return image.ToBytesFromPath(path)
}

func (a *App) cropFailed(alert string) {
	runtime.EventsEmit(a.ctx, "cropping", false)
	runtime.EventsEmit(a.ctx, "alert", alert)
}

func (a *App) RemoveBackground() []string {
	a.rmbgMu.Lock()
	path, img, model, gen := a.imgData.rembgpath, a.imgData.rembgimg, a.imgData.model, a.imgData.rembgGen
	a.imgData.rembgimg2 = nil
	a.rmbgMu.Unlock()

	if img == nil && path == "" {
		runtime.EventsEmit(a.ctx, "alert", "NO_IMAGE")
		return nil
	}

	runtime.EventsEmit(a.ctx, "removingbg", true)
	result, err := image.RemBG(model, path, img)
	if err != nil {
		runtime.EventsEmit(a.ctx, "removingbg", false)
		runtime.EventsEmit(a.ctx, "alert", err.Error())
		return nil
	}
	str, fileType, err := image.ToBase64FromBytes(result)
	if err != nil {
		runtime.EventsEmit(a.ctx, "removingbg", false)
		runtime.EventsEmit(a.ctx, "alert", err.Error())
		return nil
	}

	a.rmbgMu.Lock()
	defer a.rmbgMu.Unlock()
	if gen != a.imgData.rembgGen {
		runtime.EventsEmit(a.ctx, "removingbg", false)
		return nil
	}
	a.imgData.rembgimg2 = result
	return []string{str, fileType}
}

func (a *App) SetModel(model string) {
	a.rmbgMu.Lock()
	defer a.rmbgMu.Unlock()
	a.imgData.model = model
}

func (a *App) SaveImage(ImageType string, isJPG bool) error {
	var imgBytes []byte
	var err error

	if ImageType == "RMBG" {
		a.rmbgMu.Lock()
		imgBytes = a.imgData.rembgimg2
		a.rmbgMu.Unlock()
		isJPG = false
	} else if ImageType == "CROP" {
		a.cropMu.Lock()
		imgBytes, err = a.cropImageBytes(isJPG)
		a.cropMu.Unlock()
	}

	if err != nil {
		runtime.EventsEmit(a.ctx, "alert", err.Error())
		return err
	}
	if imgBytes == nil {
		runtime.EventsEmit(a.ctx, "alert", "NO_IMAGE")
		return nil
	}

	filePath, err := image.SaveImageDialog(a.ctx, isJPG)
	if err != nil {
		runtime.EventsEmit(a.ctx, "alert", err.Error())
		return err
	}
	if filePath == "" {
		return nil
	}

	if err := os.WriteFile(filePath, imgBytes, 0644); err != nil {
		runtime.EventsEmit(a.ctx, "alert", err.Error())
		return err
	}
	runtime.EventsEmit(a.ctx, "alert", "SAVED")
	return nil
}

func (a *App) CopyImage(ImageType string, isJPG bool) {
	var imgBytes []byte
	var err error

	if ImageType == "RMBG" {
		a.rmbgMu.Lock()
		imgBytes = a.imgData.rembgimg2
		a.rmbgMu.Unlock()
	} else if ImageType == "CROP" {
		a.cropMu.Lock()
		imgBytes, err = a.cropImageBytes(isJPG)
		a.cropMu.Unlock()
	}

	if err != nil {
		runtime.EventsEmit(a.ctx, "alert", err.Error())
		return
	}
	if imgBytes == nil {
		runtime.EventsEmit(a.ctx, "alert", "NO_IMAGE")
		return
	}

	if err := clipboard.CopyImage(imgBytes); err != nil {
		runtime.EventsEmit(a.ctx, "alert", err.Error())
		return
	}
	runtime.EventsEmit(a.ctx, "alert", "COPIED")
}

func (a *App) ClearImageMem(Imagetype string) {
	if Imagetype == "RMBG" {
		a.rmbgMu.Lock()
		a.setRMBGImage("", nil)
		a.imgData.rembgimg2 = nil
		a.rmbgMu.Unlock()
	} else if Imagetype == "CROP" {
		a.cropMu.Lock()
		a.imgData.pngcropimg = nil
		a.imgData.pngcropimgpath = ""
		a.imgData.jpgcropimg = nil
		a.imgData.jpgcropimgpath = ""
		a.cropMu.Unlock()
	} else {
		return
	}
	goruntime.GC()
}

func (a *App) CropImage(left, right, top, bottom float32, isJPG bool) []string {
	runtime.EventsEmit(a.ctx, "cropping", true)

	a.cropMu.Lock()
	defer a.cropMu.Unlock()

	imgBytes, err := a.cropImageBytes(isJPG)
	if err != nil {
		a.cropFailed(err.Error())
		return nil
	}
	if imgBytes == nil {
		a.cropFailed("NO_IMAGE")
		return nil
	}

	img, _, err := image.CropImage(imgBytes, left, right, top, bottom)
	if err != nil {
		a.cropFailed(err.Error())
		return nil
	}

	img, err = a.setCropImage("", img)
	if err != nil {
		a.cropFailed(err.Error())
		return nil
	}

	str, fileType, err := image.ToBase64FromBytes(img)
	if err != nil {
		a.cropFailed(err.Error())
		return nil
	}
	return []string{str, fileType}
}

func (a *App) FromRMBGtoCrop(t int, isJPG bool) {
	a.rmbgMu.Lock()
	imgBytes, path, output := a.imgData.rembgimg, a.imgData.rembgpath, a.imgData.rembgimg2
	a.rmbgMu.Unlock()

	a.cropMu.Lock()
	defer a.cropMu.Unlock()

	if t == 0 {
		imgBytes, err := readImage(imgBytes, path)
		if err != nil {
			runtime.EventsEmit(a.ctx, "alert", err.Error())
			return
		}
		if imgBytes != nil {
			if _, err := a.setCropImage(path, imgBytes); err != nil {
				runtime.EventsEmit(a.ctx, "alert", err.Error())
			}
		}
	} else if t == 1 {
		if output != nil {
			a.imgData.pngcropimgpath = ""
			a.imgData.pngcropimg = output
			a.imgData.jpgcropimgpath = ""
			a.imgData.jpgcropimg = nil
		}
	}
}

func (a *App) FromCroptoRMBG(isJPG bool) {
	a.cropMu.Lock()
	var imgBytes []byte
	var path string
	if isJPG {
		imgBytes, path = a.imgData.jpgcropimg, a.imgData.jpgcropimgpath
	} else {
		imgBytes, path = a.imgData.pngcropimg, a.imgData.pngcropimgpath
	}
	a.cropMu.Unlock()

	if imgBytes != nil {
		path = ""
	} else if path == "" {
		return
	}

	a.rmbgMu.Lock()
	defer a.rmbgMu.Unlock()
	a.setRMBGImage(path, imgBytes)
}

func (a *App) ToJPG() []string {
	runtime.EventsEmit(a.ctx, "cropping", true)

	a.cropMu.Lock()
	defer a.cropMu.Unlock()

	pngBytes, err := a.cropImageBytes(false)
	if err != nil {
		a.cropFailed(err.Error())
		return nil
	}
	if pngBytes != nil {
		err = image.ToJPG(pngBytes, &a.imgData.jpgcropimg)
		if err != nil {
			a.cropFailed(err.Error())
			return nil
		}
	}

	str, fileType, err := image.ToBase64FromBytes(a.imgData.jpgcropimg)
	if err != nil {
		a.cropFailed(err.Error())
		return nil
	}
	return []string{str, fileType}
}

func (a *App) ToPNG() []string {
	runtime.EventsEmit(a.ctx, "cropping", true)

	a.cropMu.Lock()
	defer a.cropMu.Unlock()

	jpgBytes, err := a.cropImageBytes(true)
	if err != nil {
		a.cropFailed(err.Error())
		return nil
	}
	if jpgBytes != nil {
		err = image.ToPNG(jpgBytes, &a.imgData.pngcropimg)
		if err != nil {
			a.cropFailed(err.Error())
			return nil
		}
	}

	str, fileType, err := image.ToBase64FromBytes(a.imgData.pngcropimg)
	if err != nil {
		a.cropFailed(err.Error())
		return nil
	}
	return []string{str, fileType}
}
