package main

import (
	"embed"
	"flag"

	"syaniclauncher/app"
	"syaniclauncher/src/helpers"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// supplied when bulding
var BUILD_ID string = "dev"

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := app.NewApp()

	// bind debug flag to DebugMode
	flag.BoolVar(&helpers.DebugMode, "debug", true, "Enable debug mode") // true for now
	flag.Parse()

	window_title := "Syanic Launcher (" + BUILD_ID + ")"
	err := wails.Run(&options.App{
		Title:  window_title,
		Width:  1300,
		Height: 750,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 10, G: 10, B: 10, A: 1}, //#0a0a0a
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
