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
	path        string
	img         []byte
	model       string
	rembgimg    []byte
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

func (a *App) HandleDrop(path string, isUrl bool) []string {
	// println("Dropped:", path, "isUrl:", isUrl)
	if isUrl {
		a.path = ""
		imgBytes, err := image.GetImageFromURL(path)
		if err != nil {
			return nil
		}
		a.img = imgBytes
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
		a.img = nil
		a.path = path
		return []string{str, fileType}
	}
}
func (a *App) OpenImage() []string {
	path, err := image.OpenImageDialog(a.ctx)
	if err != nil {
		return nil
	}
	a.img = nil
	a.path = path
	str, fileType, err := image.ToBase64FromPath(path)
	if err != nil {
		return nil
	}
	return []string{str, fileType}
}

func (a *App) RemoveBackground() []string {
	var err error
	if a.img != nil || a.path != "" {
		runtime.EventsEmit(a.ctx, "removingbg", true)
		a.rembgimg, err = image.RemBG(a.model, a.path, a.img)
		if err != nil {
			runtime.LogError(a.ctx, "Error removing background:"+err.Error())
			runtime.EventsEmit(a.ctx, "removingbg", false)
			return nil
		}
		str, fileType, err := image.ToBase64FromBytes(a.rembgimg)
		if err != nil {
			runtime.EventsEmit(a.ctx, "removingbg", false)
			return nil
		}
		// runtime.EventsEmit(a.ctx, "removingbg", false)
		return []string{str, fileType}
	}
	return nil
}

func (a *App) SetModel(model string) {
	a.model = model
}

func (a *App) SaveImage() error {
	filePath, err := image.SaveImageDialog(a.ctx)
	if err != nil {
		return err
	}
	os.WriteFile(filePath, a.rembgimg, 0644)
	return nil
}

func (a *App) CopyImage() {
	image.CopyToClipboard(a.rembgimg)
}
