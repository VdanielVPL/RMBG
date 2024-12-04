package main

import (
	"embed"
	"rmbg/utils"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist all:assets
var assets embed.FS

func main() {

	utils.InitAssets(assets)
	// Create an instance of the app structure
	app := NewApp()

	windowsOptions := &windows.Options{}

	if utils.IsWindows11OrGreater() {
		windowsOptions.WindowIsTranslucent = true
		windowsOptions.BackdropType = windows.Mica
	}

	err := wails.Run(&options.App{
		Title:     "RMBG",
		Width:     1200,
		Height:    700,
		MinWidth:  920,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: utils.BgColor(),
		OnStartup:        app.startup,
		Windows:          windowsOptions,
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop:  true,
			CSSDropProperty: "--wails-drop-target",
			CSSDropValue:    "drop",
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
