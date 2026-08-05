package main

import (
	"context"
	"embed"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "Power Plan Switcher",
		Width:     900,
		Height:    640,
		MinWidth:  700,
		MinHeight: 480,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 20, G: 22, B: 26, A: 1},
		// Closing the window hides it to the tray instead of exiting; the
		// tray's "Quit" item is what actually terminates the app.
		HideWindowOnClose: true,
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			go systray.Run(app.onTrayReady, app.onTrayExit)
		},
		OnShutdown: func(ctx context.Context) {
			systray.Quit()
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
