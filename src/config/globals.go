package config

// this file is supposed to contain default configuration values
// that gets changed on initialization of the program

// and are accessed globally like config.AppDir

var (
	BackendReady bool = false // gv for backend ready check that runs in app startup

	AppDir string // global value, set by initConfig, root of the app's data directory
)
