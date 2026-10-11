package app

import (
	"fmt"
	"syaniclauncher/src/helpers"
)

// dev: checkConnection function to check connection between backend and frontend
func (a *App) CheckConnection() string {
	helpers.LauncherLog("Frontend connection check recived")
	return fmt.Sprintln("Absolute good :)")
}

func (a *App) BrowserOpenURL(url string) {
	helpers.OpenInBrowser(a.ctx, url)
}
