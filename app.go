package main

import (
	"context"
	"os"
	"rmbg/utils"
	"rmbg/utils/image"
	"syscall"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	kernel32    *syscall.LazyDLL
	user32      *syscall.LazyDLL
	LangStrings map[string]string
	IsDarkMode  bool
	AccentColor string
	imgData     *imgDataStruct
}

type imgDataStruct struct {
	CF_PNG      uint32
	rembgpath   string
	rembgimg    []byte
	rembgimg2   []byte
	model       string
	cropimgpath string
	cropimg     []byte
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.kernel32 = syscall.NewLazyDLL("kernel32.dll")
	a.user32 = syscall.NewLazyDLL("user32.dll")

	if langData, err := utils.LoadLang(a.kernel32); err == nil {
		a.LangStrings = langData
	} else {
		runtime.LogError(ctx, "Error loading lang file:"+err.Error())
	}

	if accentColor, err := utils.GetAccentColor(); err == nil {
		a.AccentColor = accentColor
	} else {
		runtime.LogError(ctx, "Error getting accent color:"+err.Error())
	}

	a.IsDarkMode = utils.IsDarkMode()

	formatName, _ := syscall.UTF16PtrFromString("PNG")

	a.imgData = &imgDataStruct{
		model:  "u2net",
		CF_PNG: image.RegisterClipboardFormat(formatName, a.user32),
	}
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
	// println("Dropped:", path, "isUrl:", isUrl)
	if isUrl {
		if Imagetype == "RMBG" {
			a.imgData.rembgpath = ""
		} else if Imagetype == "CROP" {
			a.imgData.cropimgpath = ""
		}
		imgBytes, err := image.GetImageFromURL(path)
		if err != nil {
			return nil
		}
		if Imagetype == "RMBG" {
			a.imgData.rembgimg = imgBytes
		} else if Imagetype == "CROP" {
			a.imgData.cropimg = imgBytes
		}
		str, fileType, err := image.ToBase64FromBytes(imgBytes)
		if err != nil {
			return nil
		}
		return []string{str, fileType}
	} else {
		str, fileType, err := image.ToBase64FromPath(path)
		if err != nil {
			return nil
		}
		if Imagetype == "RMBG" {
			a.imgData.rembgimg = nil
			a.imgData.rembgpath = path
		} else if Imagetype == "CROP" {
			a.imgData.cropimg = nil
			a.imgData.cropimgpath = path
		}
		return []string{str, fileType}
	}
}
func (a *App) OpenImage(Imagetype string) []string {
	path, err := image.OpenImageDialog(a.ctx)
	if err != nil {
		return nil
	}
	if path != "" {
		if Imagetype == "RMBG" {
			a.imgData.rembgimg = nil
			a.imgData.rembgpath = path
		} else if Imagetype == "CROP" {
			a.imgData.cropimg = nil
			a.imgData.cropimgpath = path
		}
		str, fileType, err := image.ToBase64FromPath(path)
		if err != nil {
			return nil
		}
		return []string{str, fileType}
	} else {
		return nil
	}
}

func (a *App) RemoveBackground() []string {
	var err error
	if a.imgData.rembgimg != nil || a.imgData.rembgpath != "" {
		runtime.EventsEmit(a.ctx, "removingbg", true)
		a.imgData.rembgimg2, err = image.RemBG(a.imgData.model, a.imgData.rembgpath, a.imgData.rembgimg)
		if err != nil {
			runtime.EventsEmit(a.ctx, "removingbg", false)
			runtime.EventsEmit(a.ctx, "alert", err.Error())
			return nil
		}
		str, fileType, err := image.ToBase64FromBytes(a.imgData.rembgimg2)
		if err != nil {
			runtime.EventsEmit(a.ctx, "removingbg", false)
			runtime.EventsEmit(a.ctx, "alert", err.Error())
			return nil
		}
		return []string{str, fileType}
	}
	// println("No image to remove background")
	return nil
}

func (a *App) SetModel(model string) {
	a.imgData.model = model
}

func (a *App) SaveImage(ImageType string) error {
	if ImageType == "RMBG" && a.imgData.rembgimg2 != nil {
		filePath, err := image.SaveImageDialog(a.ctx)
		if err != nil {
			return err
		}
		if filePath != "" {
			os.WriteFile(filePath, a.imgData.rembgimg2, 0644)
			runtime.EventsEmit(a.ctx, "alert", "SAVED")
		}
	} else if ImageType == "CROP" && a.imgData.cropimg != nil {
		filePath, err := image.SaveImageDialog(a.ctx)
		if err != nil {
			return err
		}
		if filePath != "" {
			os.WriteFile(filePath, a.imgData.cropimg, 0644)
			runtime.EventsEmit(a.ctx, "alert", "SAVED")
		}
	} else if ImageType == "CROP" && a.imgData.cropimgpath != "" {
		filePath, err := image.SaveImageDialog(a.ctx)
		if err != nil {
			return err
		}
		if filePath != "" {
			imgBytes, err := image.ToBytesFromPath(a.imgData.cropimgpath)
			if err != nil {
				return err
			}
			os.WriteFile(filePath, imgBytes, 0644)
			runtime.EventsEmit(a.ctx, "alert", "SAVED")
		}
	}
	return nil
}

func (a *App) CopyImage(ImageType string) {
	if ImageType == "RMBG" && a.imgData.rembgimg2 != nil {
		image.CopyToClipboard(a.imgData.rembgimg2, a.kernel32, a.user32, a.imgData.CF_PNG)
	} else if ImageType == "CROP" && a.imgData.cropimg != nil {
		image.CopyToClipboard(a.imgData.cropimg, a.kernel32, a.user32, a.imgData.CF_PNG)
	} else if ImageType == "CROP" && a.imgData.cropimgpath != "" {
		imgBytes, _ := image.ToBytesFromPath(a.imgData.cropimgpath)
		image.CopyToClipboard(imgBytes, a.kernel32, a.user32, a.imgData.CF_PNG)
	} else {
		return
	}
	runtime.EventsEmit(a.ctx, "alert", "COPIED")
}

func (a *App) ClearImageMem(Imagetype string) {
	if Imagetype == "RMBG" {
		a.imgData.rembgimg = nil
		a.imgData.rembgimg2 = nil
		a.imgData.rembgpath = ""
	} else if Imagetype == "CROP" {
		a.imgData.cropimg = nil
		a.imgData.cropimgpath = ""
	}
}

func (a *App) CropImage(left, right, top, bottom float32) []string {
	runtime.EventsEmit(a.ctx, "cropping", true)
	if a.imgData.cropimg == nil {
		if a.imgData.cropimgpath != "" {
			imgBytes, err := image.ToBytesFromPath(a.imgData.cropimgpath)
			if err != nil {
				runtime.EventsEmit(a.ctx, "cropping", false)
				return nil
			}
			a.imgData.cropimgpath = ""
			a.imgData.cropimg = imgBytes
		} else {
			runtime.EventsEmit(a.ctx, "cropping", false)
			return nil
		}
	}
	img := image.CropImage(a.imgData.cropimg, left, right, top, bottom)
	if img != nil {
		a.imgData.cropimg = img
		str, fileType, err := image.ToBase64FromBytes(img)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			return nil
		}
		return []string{str, fileType}
	} else {
		runtime.EventsEmit(a.ctx, "cropping", false)
		return nil
	}
}

func (a *App) FromRMBGtoCrop(t int) {
	if t == 0 {
		if a.imgData.rembgimg != nil {
			a.imgData.cropimgpath = ""
			a.imgData.cropimg = a.imgData.rembgimg
		} else if a.imgData.rembgpath != "" {
			imgBytes, _ := image.ToBytesFromPath(a.imgData.rembgpath)
			a.imgData.cropimgpath = ""
			a.imgData.cropimg = imgBytes
		}
	} else if t == 1 {
		if a.imgData.rembgimg2 != nil {
			a.imgData.cropimgpath = ""
			a.imgData.cropimg = a.imgData.rembgimg2
		}
	}
}

func (a *App) FromCroptoRMBG() {
	if a.imgData.cropimg != nil {
		a.imgData.rembgpath = ""
		a.imgData.rembgimg = a.imgData.cropimg
	} else if a.imgData.cropimgpath != "" {
		a.imgData.rembgimg = nil
		a.imgData.rembgpath = a.imgData.cropimgpath
	}
}
