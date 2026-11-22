package helpers

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// show error dialog
func ShowError(appContext context.Context, message string) {
	runtime.MessageDialog(appContext, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   "Error",
		Message: message,
	})
}

// opens a link in browser
func OpenInBrowser(appContext context.Context, url string) {
	runtime.BrowserOpenURL(appContext, url)
}
