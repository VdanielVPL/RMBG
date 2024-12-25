package main

import (
	"context"
	"os"
	"rmbg/utils"
	"rmbg/utils/image"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	LangStrings map[string]string
	IsDarkMode  bool
	AccentColor string
	model       string
	rembgpath   string
	rembgimg    []byte
	rembgimg2   []byte
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
	if langData, err := utils.LoadLang(); err == nil {
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
	a.model = "u2net"
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
			a.rembgpath = ""
		} else if Imagetype == "CROP" {
			a.cropimgpath = ""
		}
		imgBytes, err := image.GetImageFromURL(path)
		if err != nil {
			return nil
		}
		if Imagetype == "RMBG" {
			a.rembgimg = imgBytes
		} else if Imagetype == "CROP" {
			a.cropimg = imgBytes
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
			a.rembgimg = nil
			a.rembgpath = path
		} else if Imagetype == "CROP" {
			a.cropimg = nil
			a.cropimgpath = path
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
			a.rembgimg = nil
			a.rembgpath = path
		} else if Imagetype == "CROP" {
			a.cropimg = nil
			a.cropimgpath = path
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
	if a.rembgimg != nil || a.rembgpath != "" {
		runtime.EventsEmit(a.ctx, "removingbg", true)
		a.rembgimg2, err = image.RemBG(a.model, a.rembgpath, a.rembgimg)
		if err != nil {
			runtime.EventsEmit(a.ctx, "removingbg", false)
			runtime.EventsEmit(a.ctx, "alert", err.Error())
			return nil
		}
		str, fileType, err := image.ToBase64FromBytes(a.rembgimg2)
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
	a.model = model
}

func (a *App) SaveImage(ImageType string) error {
	if ImageType == "RMBG" && a.rembgimg2 != nil {
		filePath, err := image.SaveImageDialog(a.ctx)
		if err != nil {
			return err
		}
		if filePath != "" {
			os.WriteFile(filePath, a.rembgimg2, 0644)
			runtime.EventsEmit(a.ctx, "alert", "SAVED")
		}
	} else if ImageType == "CROP" && a.cropimg != nil {
		filePath, err := image.SaveImageDialog(a.ctx)
		if err != nil {
			return err
		}
		if filePath != "" {
			os.WriteFile(filePath, a.cropimg, 0644)
			runtime.EventsEmit(a.ctx, "alert", "SAVED")
		}
	} else if ImageType == "CROP" && a.cropimgpath != "" {
		filePath, err := image.SaveImageDialog(a.ctx)
		if err != nil {
			return err
		}
		if filePath != "" {
			imgBytes, err := image.ToBytesFromPath(a.cropimgpath)
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
	if ImageType == "RMBG" && a.rembgimg2 != nil {
		image.CopyToClipboard(a.rembgimg2)
	} else if ImageType == "CROP" && a.cropimg != nil {
		image.CopyToClipboard(a.cropimg)
	} else if ImageType == "CROP" && a.cropimgpath != "" {
		imgBytes, _ := image.ToBytesFromPath(a.cropimgpath)
		image.CopyToClipboard(imgBytes)
	}
	runtime.EventsEmit(a.ctx, "alert", "COPIED")
}

func (a *App) ClearImageMem(Imagetype string) {
	if Imagetype == "RMBG" {
		a.rembgimg = nil
		a.rembgimg2 = nil
		a.rembgpath = ""
	} else if Imagetype == "CROP" {
		a.cropimg = nil
		a.cropimgpath = ""
	}
}

func (a *App) CropImage(left, right, top, bottom float32) []string {
	runtime.EventsEmit(a.ctx, "cropping", true)
	if a.cropimg == nil {
		if a.cropimgpath != "" {
			imgBytes, err := image.ToBytesFromPath(a.cropimgpath)
			if err != nil {
				runtime.EventsEmit(a.ctx, "cropping", false)
				return nil
			}
			a.cropimgpath = ""
			a.cropimg = imgBytes
		} else {
			runtime.EventsEmit(a.ctx, "cropping", false)
			return nil
		}
	}
	img := image.CropImage(a.cropimg, left, right, top, bottom)
	if img != nil {
		a.cropimg = img
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
		if a.rembgimg != nil {
			a.cropimgpath = ""
			a.cropimg = a.rembgimg
		} else if a.rembgpath != "" {
			imgBytes, _ := image.ToBytesFromPath(a.rembgpath)
			a.cropimgpath = ""
			a.cropimg = imgBytes
		}
	} else if t == 1 {
		if a.rembgimg2 != nil {
			a.cropimgpath = ""
			a.cropimg = a.rembgimg2
		}
	}
}

func (a *App) FromCroptoRMBG() {
	if a.cropimg != nil {
		a.rembgpath = ""
		a.rembgimg = a.cropimg
	} else if a.cropimgpath != "" {
		a.rembgimg = nil
		a.rembgpath = a.cropimgpath
	}
}
