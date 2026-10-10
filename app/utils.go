package app

import (
	"fmt"
	"syaniclauncher/src/helpers"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// dev: checkConnection function to check connection between backend and frontend
func (a *App) CheckConnection() string {
	helpers.LauncherLog("Frontend connection check recived")
	return fmt.Sprintln("Absolute good :)")
}

func (a *App) BrowserOpenURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}
