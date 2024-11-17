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
	a.IsDarkMode = utils.IsDarkMode()
	a.model = "u2net"
}

func (a *App) GetLangStrings() map[string]string {
	return a.LangStrings
}

func (a *App) GetDarkMode() bool {
	return a.IsDarkMode
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
			runtime.LogError(a.ctx, "Error removing background:"+err.Error())
			runtime.EventsEmit(a.ctx, "removingbg", false)
			return nil
		}
		str, fileType, err := image.ToBase64FromBytes(a.rembgimg2)
		if err != nil {
			runtime.EventsEmit(a.ctx, "removingbg", false)
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

func (a *App) SaveImage() error {
	if a.rembgimg2 != nil {
		filePath, err := image.SaveImageDialog(a.ctx)
		if err != nil {
			return err
		}
		if filePath != "" {
			os.WriteFile(filePath, a.rembgimg2, 0644)
		}
		return nil
	}
	return nil
}

func (a *App) CopyImage() {
	if a.rembgimg2 != nil {
		image.CopyToClipboard(a.rembgimg2)
	}
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
