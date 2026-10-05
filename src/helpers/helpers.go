package helpers

import (
	"context"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// make ctx available globally
var appContext context.Context

func SetContext(ctx context.Context) {
	appContext = ctx
}

func LauncherLog(msg string) {
	// todo: add debug check
	log.Println(msg)
}

func WarnLauncherLog(msg string) {
	// todo: add debug check
	log.Println()
}

func FatalLauncherLog(err error) {
	// todo: add debug check
	log.Fatal(err)
}

// error dialog
func ShowError(message string) {
	runtime.MessageDialog(appContext, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   "Error",
		Message: message,
	})
}
