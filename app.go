package main

import (
	"context"
	"net/http"
	"os"
	"rmbg/utils"
	"rmbg/utils/image"
	"rmbg/utils/image/clipboard"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// App struct
type App struct {
	app         *application.App
	LangStrings map[string]string
	IsDarkMode  bool
	AccentColor string
	imgData     *imgDataStruct
	lastImageID atomic.Uint64

	rmbgMu    sync.Mutex // not held while rembg runs
	cropMu    sync.Mutex // never held together with rmbgMu
	rembgRuns sync.WaitGroup
}

// Every image has an ID that is in its URL, an ID always shows the same image
// and stays when the image moves between the views, so the frontend has it in cache.
// IDs start from the startup time, so the cache never shows an image from an earlier run
type imgDataStruct struct {
	rembgpath      string
	rembgimg       []byte
	rembgID        uint64
	rembgimg2      []byte
	rembg2ID       uint64
	rembgGen       int                // changes with every new RMBG input
	rembgCancel    context.CancelFunc // kills the running rembg
	rmbgLoad       int                // changes with every image loaded into the tab,
	cropLoad       int                // an older load that finishes later is dropped
	model          string
	pngcropimgpath string
	pngcropimg     []byte
	pngcropID      uint64
	jpgcropimgpath string
	jpgcropimg     []byte
	jpgcropID      uint64
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// ServiceStartup is called by Wails when the app starts
func (a *App) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	a.app = application.Get()

	if langData, err := utils.LoadLang(); err == nil {
		a.LangStrings = langData
	} else {
		a.app.Logger.Error(err.Error())
	}

	if accentColor, err := utils.GetAccentColor(); err == nil {
		a.AccentColor = accentColor
	} else {
		a.app.Logger.Error(err.Error())
	}

	a.IsDarkMode = utils.IsDarkMode()

	a.imgData = &imgDataStruct{
		model: "u2net",
	}
	a.lastImageID.Store(uint64(time.Now().UnixNano()))

	if err := clipboard.Init(); err != nil {
		a.app.Logger.Error(err.Error())
	}
	return nil
}

// ServiceShutdown is called by Wails when the app is closing
func (a *App) ServiceShutdown() error {
	a.rmbgMu.Lock()
	a.stopRembg()
	a.rmbgMu.Unlock()
	a.rembgRuns.Wait()

	clipboard.Close()
	return nil
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

func (a *App) alert(key string) {
	a.app.Event.Emit("alert", key)
}

// alertError shows a translated message, technical details only go to the log
func (a *App) alertError(key string, err error) {
	if err != nil {
		if msg := err.Error(); msg == "NO_IMAGE" || msg == "INVALID_IMAGE" || msg == "REMBG_NOT_FOUND" {
			key = msg
		} else {
			a.app.Logger.Error(msg)
		}
	}
	a.alert(key)
}

func (a *App) cropFailed(key string, err error) {
	a.app.Event.Emit("cropping", false)
	a.alertError(key, err)
}

func (a *App) newImageID() uint64 {
	return a.lastImageID.Add(1)
}

// toFrontend returns the image as [URL, fileType]
func toFrontend(id uint64, imgBytes []byte) []string {
	return []string{"/images/" + strconv.FormatUint(id, 10), http.DetectContentType(imgBytes)}
}

// ServeHTTP serves the images to the frontend, Wails mounts it at /images
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(strings.TrimPrefix(r.URL.Path, "/"), 10, 64)
	var imgBytes []byte
	if err == nil && id != 0 {
		imgBytes, err = a.imageByID(id)
	}
	if err != nil || imgBytes == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", http.DetectContentType(imgBytes))
	w.Header().Set("Cache-Control", "max-age=31536000, immutable")
	w.Write(imgBytes)
}

func (a *App) imageByID(id uint64) ([]byte, error) {
	a.rmbgMu.Lock()
	imgBytes, path, inputID := a.imgData.rembgimg, a.imgData.rembgpath, a.imgData.rembgID
	output, outputID := a.imgData.rembgimg2, a.imgData.rembg2ID
	a.rmbgMu.Unlock()

	switch id {
	case inputID:
		return readImage(imgBytes, path)
	case outputID:
		return output, nil
	}

	a.cropMu.Lock()
	defer a.cropMu.Unlock()
	switch id {
	case a.imgData.pngcropID:
		return a.cropImageBytes(false)
	case a.imgData.jpgcropID:
		return a.cropImageBytes(true)
	}
	return nil, nil
}

func (a *App) HandleDrop(ImageType string, path string, isUrl bool) []string {
	load := a.newLoad(ImageType)
	var imgBytes []byte
	var err error

	if isUrl {
		imgBytes, err = image.GetImageFromURL(path)
		path = ""
	} else {
		imgBytes, err = image.ToBytesFromPath(path)
	}
	return a.setInputImage(ImageType, load, path, imgBytes, err)
}

func (a *App) OpenImage(ImageType string) []string {
	path := image.OpenImageDialog()
	if path == "" {
		return nil
	}

	load := a.newLoad(ImageType)
	imgBytes, err := image.ToBytesFromPath(path)
	return a.setInputImage(ImageType, load, path, imgBytes, err)
}

// filesDropped passes the file dropped on an image container (its id is the tab) to the frontend,
// which loads it with HandleDrop. Images dragged from a browser have no path, they are loaded from the URL
func (a *App) filesDropped(event *application.WindowEvent) {
	files, target := event.Context().DroppedFiles(), event.Context().DropTargetDetails()
	if len(files) == 0 || files[0] == "" || target == nil {
		return
	}
	a.app.Event.Emit("filedrop", []string{target.ElementID, files[0]})
}

func (a *App) newLoad(ImageType string) int {
	switch ImageType {
	case "RMBG":
		a.rmbgMu.Lock()
		defer a.rmbgMu.Unlock()
		a.imgData.rmbgLoad++
		return a.imgData.rmbgLoad
	case "CROP":
		a.cropMu.Lock()
		defer a.cropMu.Unlock()
		a.imgData.cropLoad++
		return a.imgData.cropLoad
	}
	return 0
}

func (a *App) setInputImage(ImageType string, load int, path string, imgBytes []byte, err error) []string {
	switch ImageType {
	case "RMBG":
		a.rmbgMu.Lock()
		defer a.rmbgMu.Unlock()
		if load != a.imgData.rmbgLoad {
			return nil
		}
	case "CROP":
		a.cropMu.Lock()
		defer a.cropMu.Unlock()
		if load != a.imgData.cropLoad {
			return nil
		}
	default:
		return nil
	}

	if err != nil {
		a.alertError("LOAD_FAILED", err)
		return nil
	}
	if !image.IsSupportedImage(imgBytes) {
		a.alert("INVALID_IMAGE")
		return nil
	}

	if ImageType == "RMBG" {
		a.setRMBGImage(path, imgBytes, a.newImageID())
		return toFrontend(a.imgData.rembgID, imgBytes)
	}

	result, err := a.setCropImage(path, imgBytes, a.newImageID())
	if err != nil {
		a.alertError("LOAD_FAILED", err)
		return nil
	}
	return result
}

func (a *App) setRMBGImage(path string, imgBytes []byte, id uint64) {
	a.imgData.rembgpath = path
	a.imgData.rembgimg = nil
	if path == "" {
		a.imgData.rembgimg = imgBytes
	}
	a.imgData.rembgID = id
	a.imgData.rembgimg2, a.imgData.rembg2ID = nil, 0
	a.imgData.rembgGen++
	a.stopRembg()
}

// stopRembg kills the running rembg, its result would be dropped anyway
func (a *App) stopRembg() {
	if a.imgData.rembgCancel != nil {
		a.imgData.rembgCancel()
		a.imgData.rembgCancel = nil
	}
}

func (a *App) RemoveBackground() []string {
	a.rmbgMu.Lock()
	path, img, model, gen := a.imgData.rembgpath, a.imgData.rembgimg, a.imgData.model, a.imgData.rembgGen
	a.imgData.rembgimg2, a.imgData.rembg2ID = nil, 0
	if img == nil && path == "" {
		a.rmbgMu.Unlock()
		a.alert("NO_IMAGE")
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a.imgData.rembgCancel = cancel
	a.rembgRuns.Add(1)
	defer a.rembgRuns.Done()
	a.rmbgMu.Unlock()

	a.app.Event.Emit("removingbg", true)
	output, err := image.RemBG(ctx, model, path, img)
	if err != nil {
		a.app.Event.Emit("removingbg", false)
		if ctx.Err() == nil { // canceled because the image changed, not an error
			a.alertError("REMBG_FAILED", err)
		}
		return nil
	}

	a.rmbgMu.Lock()
	defer a.rmbgMu.Unlock()
	if gen != a.imgData.rembgGen {
		a.app.Event.Emit("removingbg", false)
		return nil
	}
	a.imgData.rembgimg2, a.imgData.rembg2ID = output, a.newImageID()
	return toFrontend(a.imgData.rembg2ID, output)
}

func (a *App) SetModel(model string) {
	a.rmbgMu.Lock()
	defer a.rmbgMu.Unlock()
	a.imgData.model = model
}

// JPG and PNG files are kept as paths, other formats are converted to PNG and kept in memory
func (a *App) setCropImage(path string, imgBytes []byte, id uint64) ([]string, error) {
	imgBytes, isJPG, converted, err := image.ToSupportedFormat(imgBytes)
	if err != nil {
		return nil, err
	}
	if converted {
		path, id = "", a.newImageID()
	}

	var inMemory []byte
	if path == "" {
		inMemory = imgBytes
	}

	a.clearCropImage()
	if isJPG {
		a.imgData.jpgcropimgpath, a.imgData.jpgcropimg, a.imgData.jpgcropID = path, inMemory, id
	} else {
		a.imgData.pngcropimgpath, a.imgData.pngcropimg, a.imgData.pngcropID = path, inMemory, id
	}
	return toFrontend(id, imgBytes), nil
}

func (a *App) clearCropImage() {
	a.imgData.pngcropimgpath, a.imgData.pngcropimg, a.imgData.pngcropID = "", nil, 0
	a.imgData.jpgcropimgpath, a.imgData.jpgcropimg, a.imgData.jpgcropID = "", nil, 0
}

// cropImage returns the crop image in the given format, only one of bytes and path is set
func (a *App) cropImage(isJPG bool) ([]byte, string) {
	if isJPG {
		return a.imgData.jpgcropimg, a.imgData.jpgcropimgpath
	}
	return a.imgData.pngcropimg, a.imgData.pngcropimgpath
}

func (a *App) cropImageBytes(isJPG bool) ([]byte, error) {
	return readImage(a.cropImage(isJPG))
}

func readImage(imgBytes []byte, path string) ([]byte, error) {
	if imgBytes != nil || path == "" {
		return imgBytes, nil
	}
	return image.ToBytesFromPath(path)
}

func (a *App) CropImage(left, right, top, bottom float32, isJPG bool) []string {
	a.app.Event.Emit("cropping", true)

	a.cropMu.Lock()
	defer a.cropMu.Unlock()

	imgBytes, err := a.cropImageBytes(isJPG)
	if err != nil {
		a.cropFailed("CROP_FAILED", err)
		return nil
	}
	if imgBytes == nil {
		a.cropFailed("NO_IMAGE", nil)
		return nil
	}

	img, err := image.CropImage(imgBytes, left, right, top, bottom)
	if err != nil {
		a.cropFailed("CROP_FAILED", err)
		return nil
	}

	result, err := a.setCropImage("", img, a.newImageID())
	if err != nil {
		a.cropFailed("CROP_FAILED", err)
		return nil
	}
	return result
}

func (a *App) ToJPG() []string {
	return a.convertCropImage(true)
}

func (a *App) ToPNG() []string {
	return a.convertCropImage(false)
}

// convertCropImage makes the crop image in the other format and keeps both,
// when there is nothing to convert it returns the one that is already there
func (a *App) convertCropImage(toJPG bool) []string {
	a.app.Event.Emit("cropping", true)

	a.cropMu.Lock()
	defer a.cropMu.Unlock()

	convert, target, targetID := image.ToPNG, &a.imgData.pngcropimg, &a.imgData.pngcropID
	if toJPG {
		convert, target, targetID = image.ToJPG, &a.imgData.jpgcropimg, &a.imgData.jpgcropID
	}

	imgBytes, err := a.cropImageBytes(!toJPG)
	if err != nil {
		a.cropFailed("CONVERT_FAILED", err)
		return nil
	}
	if imgBytes != nil {
		converted, err := convert(imgBytes)
		if err != nil {
			a.cropFailed("CONVERT_FAILED", err)
			return nil
		}
		*target, *targetID = converted, a.newImageID()
	}

	if *target == nil {
		a.cropFailed("NO_IMAGE", nil)
		return nil
	}
	return toFrontend(*targetID, *target)
}

// outputImage returns the image that is saved or copied from the tab
func (a *App) outputImage(ImageType string, isJPG bool) ([]byte, error) {
	switch ImageType {
	case "RMBG":
		a.rmbgMu.Lock()
		defer a.rmbgMu.Unlock()
		return a.imgData.rembgimg2, nil
	case "CROP":
		a.cropMu.Lock()
		defer a.cropMu.Unlock()
		return a.cropImageBytes(isJPG)
	}
	return nil, nil
}

func (a *App) SaveImage(ImageType string, isJPG bool) {
	if ImageType == "RMBG" {
		isJPG = false
	}

	imgBytes, err := a.outputImage(ImageType, isJPG)
	if err != nil {
		a.alertError("SAVE_FAILED", err)
		return
	}
	if imgBytes == nil {
		a.alert("NO_IMAGE")
		return
	}

	filePath := image.SaveImageDialog(isJPG)
	if filePath == "" {
		return
	}

	if err := os.WriteFile(filePath, imgBytes, 0644); err != nil {
		a.alertError("SAVE_FAILED", err)
		return
	}
	a.alert("SAVED")
}

func (a *App) CopyImage(ImageType string, isJPG bool) {
	imgBytes, err := a.outputImage(ImageType, isJPG)
	if err != nil {
		a.alertError("COPY_FAILED", err)
		return
	}
	if imgBytes == nil {
		a.alert("NO_IMAGE")
		return
	}

	if err := clipboard.CopyImage(imgBytes); err != nil {
		a.alertError("COPY_FAILED", err)
		return
	}
	a.alert("COPIED")
}

func (a *App) ClearImageMem(ImageType string) {
	switch ImageType {
	case "RMBG":
		a.rmbgMu.Lock()
		a.setRMBGImage("", nil, 0)
		a.imgData.rmbgLoad++
		a.rmbgMu.Unlock()
	case "CROP":
		a.cropMu.Lock()
		a.clearCropImage()
		a.imgData.cropLoad++
		a.cropMu.Unlock()
	default:
		return
	}
	debug.FreeOSMemory()
}

// FromRMBGtoCrop moves the input (0) or the image without background (1) to the crop view,
// when nothing could be moved the crop view is cleared
func (a *App) FromRMBGtoCrop(t int) []string {
	a.rmbgMu.Lock()
	imgBytes, path, id := a.imgData.rembgimg, a.imgData.rembgpath, a.imgData.rembgID
	output, outputID := a.imgData.rembgimg2, a.imgData.rembg2ID
	a.rmbgMu.Unlock()

	a.cropMu.Lock()
	defer a.cropMu.Unlock()
	a.imgData.cropLoad++
	a.clearCropImage()

	switch t {
	case 0: // input image
		if imgBytes == nil && path == "" {
			return nil
		}
		imgBytes, err := readImage(imgBytes, path)
		var result []string
		if err == nil {
			result, err = a.setCropImage(path, imgBytes, id)
		}
		if err != nil {
			a.alertError("LOAD_FAILED", err)
			return nil
		}
		return result
	case 1: // image without background
		if output == nil {
			return nil
		}
		a.imgData.pngcropimg, a.imgData.pngcropID = output, outputID
		return toFrontend(outputID, output)
	}
	return nil
}

func (a *App) FromCroptoRMBG(isJPG bool) {
	a.cropMu.Lock()
	imgBytes, path := a.cropImage(isJPG)
	id := a.imgData.pngcropID
	if isJPG {
		id = a.imgData.jpgcropID
	}
	a.cropMu.Unlock()

	if imgBytes == nil && path == "" {
		return
	}
	if imgBytes != nil {
		path = ""
	}

	a.rmbgMu.Lock()
	defer a.rmbgMu.Unlock()
	a.setRMBGImage(path, imgBytes, id)
	a.imgData.rmbgLoad++
}
