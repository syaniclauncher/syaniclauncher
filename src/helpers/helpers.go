package helpers

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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

func OpenInBrowser(appContext context.Context, url string) {
	runtime.BrowserOpenURL(appContext, url)
}

func OfflineUUID(username string) string {
	sum := md5.Sum([]byte("OfflinePlayer:" + username))
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", sum[0:4], sum[4:6], sum[6:8], sum[8:10], sum[10:16])
}
