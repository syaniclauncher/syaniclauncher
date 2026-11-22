package helpers

import "log"

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
