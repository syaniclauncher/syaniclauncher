package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syaniclauncher/src/helpers"
)

// function to initialize %appdata%/.syaniclauncher/config.json

// if not exists, create it with default values
// if corrupted, reset it to default values
// if it exists, import the values

type Config struct {
	AppDirectory string `json:"app_dir"`
}

func DefaultConfig() Config {
	return Config{
		AppDirectory: "", //empty = default %appdata%/.syaniclauncher
	}
}

// write saves config to path, creating the file if needed
func write(path string, cfg Config) error {
	data, err := json.Marshal(cfg)

	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// quarantine moves an unusable config file aside so it can be inspected,
// then writes defaults in its place
func quarantine(configPath string) {
	backupPath := configPath + ".corrupt"

	if err := os.Rename(configPath, backupPath); err != nil {
		log.Println("[WARN] config: could not keep corrupted file:", err)
		return
	}

	log.Println("[WARN] config: corrupted, previous file kept at", backupPath)
}

// usableAppDir reports whether dir is an absolute path we can actually use & creating it if missing
func usableAppDir(dir string) bool {
	if !filepath.IsAbs(dir) {
		return false
	}

	return os.MkdirAll(dir, 0755) == nil
}

func InitConfig() (*Config, error) {
	appData, err := os.UserConfigDir()

	if err != nil {
		helpers.ShowError("Unable to get APPDATA directory.")
		return nil, err
	}

	defaultAppDir := filepath.Join(appData, ".syaniclauncher")
	// config location is fixed under default directory btw (%appdata%/.syaniclauncher/config.json)
	// even if use changes app_dir, config will always be here
	configPath := filepath.Join(defaultAppDir, "config.json")

	// create directory if missing
	if err = os.MkdirAll(defaultAppDir, 0755); err != nil {
		helpers.ShowError("Unable to create default app directory.")
		return nil, err
	}

	config := DefaultConfig() //start with default
	configData, err := os.ReadFile(configPath)

	if os.IsNotExist(err) {
		// Config doesn't exist, keep defaults.
	} else if err != nil {
		// Config exists but can't be read.
		quarantine(configPath)
	} else if err := json.Unmarshal(configData, &config); err != nil { //json.unmarshal here reads config data and writes it into out config struct
		// Config exists but contains invalid JSON.
		quarantine(configPath)
		config = DefaultConfig()
	}

	// if saved app_dir can't be used, fall back to the default app dir,
	// leaving every other field untouched & the json itself is valid,
	// so just rewrite it instead of quarantining

	// note: paths are saved with double slashes like ("G:\\test"), otherwise it decodes
	// to something else (eg. "\t" is a tab)
	if config.AppDirectory != "" && !usableAppDir(config.AppDirectory) {
		helpers.WarnLauncherLog(fmt.Sprintf("[WARN] [CONFIG-INIT] unusable app_dir %q, falling back to default", config.AppDirectory))
		config.AppDirectory = ""
	}

	// finally write the config after all this
	if err := write(configPath, config); err != nil {
		helpers.ShowError("Unable to update default config.")
		helpers.FatalLauncherLog(err)
		return nil, err
	}

	// global: declare AppDir here for launcher to use
	if config.AppDirectory == "" {
		AppDir = defaultAppDir
	} else {
		AppDir = config.AppDirectory
	}

	helpers.LauncherLog("[DEBUG] AppDir is set to: " + AppDir)
	return &config, nil
}
