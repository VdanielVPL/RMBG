package main

import (
	"embed"
	"os"
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

	if len(os.Args) > 1 {
		utils.Cli(os.Args)
	} else {
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
			DisableResize:    true,
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
}
