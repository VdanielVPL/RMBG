package main

import (
	"context"
	"os"
	"rmbg/utils"
	"rmbg/utils/image"
	"rmbg/utils/image/clipboard"
	goruntime "runtime"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	LangStrings map[string]string
	IsDarkMode  bool
	AccentColor string
	imgData     *imgDataStruct
}

type imgDataStruct struct {
	rembgpath      string
	rembgimg       []byte
	rembgimg2      []byte
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
		a.imgData.rembgpath = path
		a.imgData.rembgimg = nil
		if path == "" {
			a.imgData.rembgimg = imgBytes
		}
	} else if Imagetype == "CROP" {
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
		if a.imgData.jpgcropimg != nil {
			return a.imgData.jpgcropimg, nil
		} else if a.imgData.jpgcropimgpath != "" {
			return image.ToBytesFromPath(a.imgData.jpgcropimgpath)
		}
	} else {
		if a.imgData.pngcropimg != nil {
			return a.imgData.pngcropimg, nil
		} else if a.imgData.pngcropimgpath != "" {
			return image.ToBytesFromPath(a.imgData.pngcropimgpath)
		}
	}
	return nil, nil
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
	runtime.EventsEmit(a.ctx, "alert", "NO_IMAGE")
	return nil
}

func (a *App) SetModel(model string) {
	a.imgData.model = model
}

func (a *App) SaveImage(ImageType string, isJPG bool) error {
	var imgBytes []byte
	var err error

	if ImageType == "RMBG" {
		imgBytes = a.imgData.rembgimg2
		isJPG = false
	} else if ImageType == "CROP" {
		imgBytes, err = a.cropImageBytes(isJPG)
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
		imgBytes = a.imgData.rembgimg2
	} else if ImageType == "CROP" {
		imgBytes, err = a.cropImageBytes(isJPG)
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
		a.imgData.rembgimg = nil
		a.imgData.rembgimg2 = nil
		a.imgData.rembgpath = ""
	} else if Imagetype == "CROP" {
		a.imgData.pngcropimg = nil
		a.imgData.pngcropimgpath = ""
		a.imgData.jpgcropimg = nil
		a.imgData.jpgcropimgpath = ""
	} else {
		return
	}
	goruntime.GC()
}

func (a *App) CropImage(left, right, top, bottom float32, isJPG bool) []string {
	runtime.EventsEmit(a.ctx, "cropping", true)

	var imgByte []byte

	if isJPG {
		if a.imgData.jpgcropimg != nil {
			imgByte = a.imgData.jpgcropimg
		} else if a.imgData.jpgcropimgpath != "" {
			var err error
			imgByte, err = image.ToBytesFromPath(a.imgData.jpgcropimgpath)
			if err != nil {
				runtime.EventsEmit(a.ctx, "cropping", false)
				runtime.EventsEmit(a.ctx, "alert", err.Error())
				return nil
			}
		}
	} else {
		if a.imgData.pngcropimg != nil {
			imgByte = a.imgData.pngcropimg
		} else if a.imgData.pngcropimgpath != "" {
			var err error
			imgByte, err = image.ToBytesFromPath(a.imgData.pngcropimgpath)
			if err != nil {
				runtime.EventsEmit(a.ctx, "cropping", false)
				runtime.EventsEmit(a.ctx, "alert", err.Error())
				return nil
			}
		}
	}

	img, format, err := image.CropImage(imgByte, left, right, top, bottom)
	if err != nil {
		runtime.EventsEmit(a.ctx, "cropping", false)
		runtime.EventsEmit(a.ctx, "alert", err.Error())
		return nil
	}
	if img != nil {
		if format == "jpeg" {
			a.imgData.jpgcropimg = img
		} else {
			a.imgData.pngcropimg = img
		}
		str, fileType, err := image.ToBase64FromBytes(img)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			runtime.EventsEmit(a.ctx, "alert", err.Error())
			return nil
		}
		return []string{str, fileType}
	} else {
		runtime.EventsEmit(a.ctx, "cropping", false)
		runtime.EventsEmit(a.ctx, "alert", "NO_IMAGE")
		return nil
	}
}

func (a *App) FromRMBGtoCrop(t int, isJPG bool) {
	if t == 0 {
		imgBytes := a.imgData.rembgimg
		if imgBytes == nil && a.imgData.rembgpath != "" {
			var err error
			imgBytes, err = image.ToBytesFromPath(a.imgData.rembgpath)
			if err != nil {
				runtime.EventsEmit(a.ctx, "alert", err.Error())
				return
			}
		}
		if imgBytes != nil {
			if _, err := a.setCropImage(a.imgData.rembgpath, imgBytes); err != nil {
				runtime.EventsEmit(a.ctx, "alert", err.Error())
			}
		}
	} else if t == 1 {
		if a.imgData.rembgimg2 != nil {
			a.imgData.pngcropimgpath = ""
			a.imgData.pngcropimg = a.imgData.rembgimg2
			a.imgData.jpgcropimgpath = ""
			a.imgData.jpgcropimg = nil
		}
	}
}

func (a *App) FromCroptoRMBG(isJPG bool) {
	if isJPG {
		if a.imgData.jpgcropimg != nil {
			a.imgData.rembgpath = ""
			a.imgData.rembgimg = a.imgData.jpgcropimg
		} else if a.imgData.jpgcropimgpath != "" {
			a.imgData.rembgimg = nil
			a.imgData.rembgpath = a.imgData.jpgcropimgpath
		}
	} else {
		if a.imgData.pngcropimg != nil {
			a.imgData.rembgpath = ""
			a.imgData.rembgimg = a.imgData.pngcropimg
		} else if a.imgData.pngcropimgpath != "" {
			a.imgData.rembgimg = nil
			a.imgData.rembgpath = a.imgData.pngcropimgpath
		}
	}
}

func (a *App) ToJPG() []string {
	runtime.EventsEmit(a.ctx, "cropping", true)
	if a.imgData.pngcropimg != nil {
		err := image.ToJPG(a.imgData.pngcropimg, &a.imgData.jpgcropimg)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			runtime.EventsEmit(a.ctx, "alert", err.Error())
			return nil
		}
	} else if a.imgData.pngcropimgpath != "" {
		imgBytes, err := image.ToBytesFromPath(a.imgData.pngcropimgpath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			runtime.EventsEmit(a.ctx, "alert", err.Error())
			return nil
		}
		err = image.ToJPG(imgBytes, &a.imgData.jpgcropimg)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			runtime.EventsEmit(a.ctx, "alert", err.Error())
			return nil
		}
	}
	str, fileType, err := image.ToBase64FromBytes(a.imgData.jpgcropimg)
	if err != nil {
		runtime.EventsEmit(a.ctx, "cropping", false)
		runtime.EventsEmit(a.ctx, "alert", err.Error())
		return nil
	}
	return []string{str, fileType}
}

func (a *App) ToPNG() []string {
	runtime.EventsEmit(a.ctx, "cropping", true)
	if a.imgData.jpgcropimg != nil {
		err := image.ToPNG(a.imgData.jpgcropimg, &a.imgData.pngcropimg)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			runtime.EventsEmit(a.ctx, "alert", err.Error())
			return nil
		}
	} else if a.imgData.jpgcropimgpath != "" {
		imgBytes, err := image.ToBytesFromPath(a.imgData.jpgcropimgpath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			runtime.EventsEmit(a.ctx, "alert", err.Error())
			return nil
		}
		err = image.ToPNG(imgBytes, &a.imgData.pngcropimg)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			runtime.EventsEmit(a.ctx, "alert", err.Error())
			return nil
		}
	}
	str, fileType, err := image.ToBase64FromBytes(a.imgData.pngcropimg)
	if err != nil {
		runtime.EventsEmit(a.ctx, "cropping", false)
		runtime.EventsEmit(a.ctx, "alert", err.Error())
		return nil
	}
	return []string{str, fileType}
}
