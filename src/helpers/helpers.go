package helpers

import (
	"context"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// make ctx available globally
var DebugMode bool

func LauncherLog(msg string) {
	if DebugMode {
		log.Println("[INFO] " + msg)
	}
}

func WarnLauncherLog(msg string) {
	if DebugMode {
		log.Println("[WARN] " + msg)
	}
}

func FatalLauncherLog(err error) {
	if DebugMode {
		log.Fatal(err)
	}
}

// error dialog
func ShowError(appContext context.Context, message string) {
	runtime.MessageDialog(appContext, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   "Error",
		Message: message,
	})
}
