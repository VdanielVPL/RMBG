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
	CF_PNG         uint32
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
			a.imgData.pngcropimgpath = ""
			a.imgData.jpgcropimgpath = ""
		}
		imgBytes, err := image.GetImageFromURL(path)
		if err != nil {
			return nil
		}
		if Imagetype == "RMBG" {
			a.imgData.rembgimg = imgBytes
		} else if Imagetype == "CROP" {
			format := image.GetMimeType(imgBytes)
			if format == "image/jpeg" {
				a.imgData.jpgcropimg = imgBytes
			} else {
				a.imgData.pngcropimg = imgBytes
			}
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
			a.imgData.pngcropimg = nil
			a.imgData.jpgcropimg = nil
			format := image.GetFormatFromPath(path)
			if format == "jpg" || format == "jpeg" {
				a.imgData.jpgcropimgpath = path
			} else {
				a.imgData.pngcropimgpath = path
			}
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
			a.imgData.pngcropimg = nil
			a.imgData.jpgcropimg = nil
			format := image.GetFormatFromPath(path)
			if format == "jpg" || format == "jpeg" {
				a.imgData.jpgcropimgpath = path
			} else {
				a.imgData.pngcropimgpath = path
			}
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
		filePath, err := image.SaveImageDialog(a.ctx, false)
		if err != nil {
			return err
		}
		if filePath != "" {
			os.WriteFile(filePath, a.imgData.rembgimg2, 0644)
			runtime.EventsEmit(a.ctx, "alert", "SAVED")
		}
	} else if ImageType == "CROP" && a.imgData.pngcropimg != nil {
		filePath, err := image.SaveImageDialog(a.ctx, false)
		if err != nil {
			return err
		}
		if filePath != "" {
			os.WriteFile(filePath, a.imgData.pngcropimg, 0644)
			runtime.EventsEmit(a.ctx, "alert", "SAVED")
		}
	} else if ImageType == "CROP" && a.imgData.jpgcropimg != nil {
		filePath, err := image.SaveImageDialog(a.ctx, true)
		if err != nil {
			return err
		}
		if filePath != "" {
			os.WriteFile(filePath, a.imgData.jpgcropimg, 0644)
			runtime.EventsEmit(a.ctx, "alert", "SAVED")
		}
	} else if ImageType == "CROP" && a.imgData.pngcropimgpath != "" {
		filePath, err := image.SaveImageDialog(a.ctx, false)
		if err != nil {
			return err
		}
		if filePath != "" {
			imgBytes, err := image.ToBytesFromPath(a.imgData.pngcropimgpath)
			if err != nil {
				return err
			}
			os.WriteFile(filePath, imgBytes, 0644)
			runtime.EventsEmit(a.ctx, "alert", "SAVED")
		}
	} else if ImageType == "CROP" && a.imgData.jpgcropimgpath != "" {
		filePath, err := image.SaveImageDialog(a.ctx, true)
		if err != nil {
			return err
		}
		if filePath != "" {
			imgBytes, err := image.ToBytesFromPath(a.imgData.jpgcropimgpath)
			if err != nil {
				return err
			}
			os.WriteFile(filePath, imgBytes, 0644)
			runtime.EventsEmit(a.ctx, "alert", "SAVED")
		}
	}

	// } else if ImageType == "CROP" && (a.imgData.pngcropimgpath != "" || a.imgData.jpgcropimgpath != "") {
	// 	filePath, err := image.SaveImageDialog(a.ctx)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if filePath != "" {
	// 		if a.imgData.jpgcropimgpath != "" {
	// 			imgBytes, err := image.ToBytesFromPath(a.imgData.jpgcropimgpath)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			os.WriteFile(filePath, imgBytes, 0644)
	// 			runtime.EventsEmit(a.ctx, "alert", "SAVED")
	// 		} else if a.imgData.pngcropimgpath != "" {
	// 			imgBytes, err := image.ToBytesFromPath(a.imgData.pngcropimgpath)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			os.WriteFile(filePath, imgBytes, 0644)
	// 			runtime.EventsEmit(a.ctx, "alert", "SAVED")
	// 		}
	// 	}
	// }
	return nil
}

func (a *App) CopyImage(ImageType string) {
	//RMBG View
	if ImageType == "RMBG" && a.imgData.rembgimg2 != nil {
		image.CopyToClipboard(a.imgData.rembgimg2, a.kernel32, a.user32, a.imgData.CF_PNG)
		//Crop View
	} else if ImageType == "CROP" && a.imgData.pngcropimg != nil {
		image.CopyToClipboard(a.imgData.pngcropimg, a.kernel32, a.user32, a.imgData.CF_PNG)
	} else if ImageType == "CROP" && a.imgData.jpgcropimg != nil {
		image.CopyToClipboard(a.imgData.jpgcropimg, a.kernel32, a.user32, a.imgData.CF_PNG)
	} else if ImageType == "CROP" && a.imgData.pngcropimgpath != "" {
		imgBytes, _ := image.ToBytesFromPath(a.imgData.pngcropimgpath)
		image.CopyToClipboard(imgBytes, a.kernel32, a.user32, a.imgData.CF_PNG)
	} else if ImageType == "CROP" && a.imgData.jpgcropimgpath != "" {
		imgBytes, _ := image.ToBytesFromPath(a.imgData.jpgcropimgpath)
		image.CopyToClipboard(imgBytes, a.kernel32, a.user32, a.imgData.CF_PNG)
		//No image to copy
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
		a.imgData.pngcropimg = nil
		a.imgData.pngcropimgpath = ""
		a.imgData.jpgcropimg = nil
		a.imgData.jpgcropimgpath = ""
	}
}

func (a *App) CropImage(left, right, top, bottom float32) []string {
	runtime.EventsEmit(a.ctx, "cropping", true)
	// if a.imgData.pngcropimg == nil {
	// 	if a.imgData.pngcropimgpath != "" {
	// 		imgBytes, err := image.ToBytesFromPath(a.imgData.pngcropimgpath)
	// 		if err != nil {
	// 			runtime.EventsEmit(a.ctx, "cropping", false)
	// 			return nil
	// 		}
	// 		a.imgData.pngcropimgpath = ""
	// 		a.imgData.pngcropimg = imgBytes
	// 	} else {
	// 		runtime.EventsEmit(a.ctx, "cropping", false)
	// 		return nil
	// 	}
	// } else {
	// 	if a.imgData.jpgcropimgpath != "" {
	// 		imgBytes, err := image.ToBytesFromPath(a.imgData.jpgcropimgpath)
	// 		if err != nil {
	// 			runtime.EventsEmit(a.ctx, "cropping", false)
	// 			return nil
	// 		}
	// 		a.imgData.jpgcropimgpath = ""
	// 		a.imgData.jpgcropimg = imgBytes
	// 	} else {
	// 		runtime.EventsEmit(a.ctx, "cropping", false)
	// 		return nil
	// 	}
	// }
	var imgByte []byte

	if a.imgData.pngcropimg != nil {
		imgByte = a.imgData.pngcropimg
	} else if a.imgData.pngcropimgpath != "" {
		imgByte, _ = image.ToBytesFromPath(a.imgData.pngcropimgpath)
	} else if a.imgData.jpgcropimg != nil {
		imgByte = a.imgData.jpgcropimg
	} else if a.imgData.jpgcropimgpath != "" {
		imgByte, _ = image.ToBytesFromPath(a.imgData.jpgcropimgpath)
	}

	img, format := image.CropImage(imgByte, left, right, top, bottom)
	if img != nil {
		if format == "jpeg" {
			a.imgData.jpgcropimg = img
		} else {
			a.imgData.pngcropimg = img
		}
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
			a.imgData.pngcropimgpath = ""
			a.imgData.pngcropimg = a.imgData.rembgimg
			a.imgData.jpgcropimgpath = ""
			a.imgData.jpgcropimg = nil
		} else if a.imgData.rembgpath != "" {
			imgBytes, _ := image.ToBytesFromPath(a.imgData.rembgpath)
			a.imgData.pngcropimgpath = ""
			a.imgData.pngcropimg = imgBytes
			a.imgData.jpgcropimgpath = ""
			a.imgData.jpgcropimg = nil
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

func (a *App) FromCroptoRMBG() {
	if a.imgData.pngcropimg != nil {
		a.imgData.rembgpath = ""
		a.imgData.rembgimg = a.imgData.pngcropimg
	} else if a.imgData.pngcropimgpath != "" {
		a.imgData.rembgimg = nil
		a.imgData.rembgpath = a.imgData.pngcropimgpath
	} else if a.imgData.jpgcropimg != nil {
		a.imgData.rembgpath = ""
		a.imgData.rembgimg = a.imgData.jpgcropimg
	} else if a.imgData.jpgcropimgpath != "" {
		a.imgData.rembgimg = nil
		a.imgData.rembgpath = a.imgData.jpgcropimgpath
	}
}

func (a *App) ToJPG() []string {
	if a.imgData.pngcropimg != nil {
		err := image.ToJPG(a.imgData.pngcropimg, &a.imgData.jpgcropimg)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			return nil
		}
	} else if a.imgData.pngcropimgpath != "" {
		imgBytes, err := image.ToBytesFromPath(a.imgData.pngcropimgpath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			return nil
		}
		err = image.ToJPG(imgBytes, &a.imgData.jpgcropimg)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			return nil
		}
	}
	str, fileType, err := image.ToBase64FromBytes(a.imgData.jpgcropimg)
	if err != nil {
		runtime.EventsEmit(a.ctx, "cropping", false)
		return nil
	}
	return []string{str, fileType}
}

func (a *App) ToPNG() []string {
	if a.imgData.jpgcropimg != nil {
		err := image.ToPNG(a.imgData.jpgcropimg, &a.imgData.pngcropimg)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			return nil
		}
	} else if a.imgData.jpgcropimgpath != "" {
		imgBytes, err := image.ToBytesFromPath(a.imgData.jpgcropimgpath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			return nil
		}
		err = image.ToPNG(imgBytes, &a.imgData.pngcropimg)
		if err != nil {
			runtime.EventsEmit(a.ctx, "cropping", false)
			return nil
		}
	}
	str, fileType, err := image.ToBase64FromBytes(a.imgData.pngcropimg)
	if err != nil {
		runtime.EventsEmit(a.ctx, "cropping", false)
		return nil
	}
	return []string{str, fileType}
}
