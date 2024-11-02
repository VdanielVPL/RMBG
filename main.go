package main

import (
	"embed"
	"rmbg/utils"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	windowsOptions := &windows.Options{}

	if utils.IsWindows11OrGreater() {
		windowsOptions.WindowIsTranslucent = true
		windowsOptions.BackdropType = windows.Mica
	}

	err := wails.Run(&options.App{
		Title:  "rmbg",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: utils.BgColor(),
		OnStartup:        app.startup,
		Windows:          windowsOptions,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
