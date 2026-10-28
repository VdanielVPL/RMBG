package main

import (
	"embed"
	"os"
	"rmbg/utils"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist all:assets
var assets embed.FS

func main() {

	utils.InitAssets(assets)

	if len(os.Args) > 1 {
		if err := utils.Cli(os.Args); err != nil {
			println("Error:", err.Error())
		}
	} else {
		rmbg := NewApp()

		app := application.New(application.Options{
			Name: "RMBG",
			Services: []application.Service{
				application.NewService(rmbg),
			},
			Assets: application.AssetOptions{
				Handler: application.AssetFileServerFS(assets),
			},
		})

		windowOptions := application.WebviewWindowOptions{
			Title:                      "RMBG",
			Width:                      1200,
			Height:                     700,
			MinWidth:                   920,
			MinHeight:                  600,
			DisableResize:              true,
			BackgroundColour:           utils.BgColor(),
			EnableFileDrop:             true,
			DefaultContextMenuDisabled: true,
		}

		if utils.IsWindows11OrGreater() {
			windowOptions.BackgroundType = application.BackgroundTypeTranslucent
			windowOptions.Windows.BackdropType = application.Mica
		}

		window := app.Window.NewWithOptions(windowOptions)
		window.OnWindowEvent(events.Common.WindowFilesDropped, rmbg.filesDropped)

		if err := app.Run(); err != nil {
			println("Error:", err.Error())
		}
	}
}
