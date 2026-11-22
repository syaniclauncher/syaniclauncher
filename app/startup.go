package app

import (
	"context"
	"syaniclauncher/src/helpers"
	"syaniclauncher/src/launcher_accounts"
	"syaniclauncher/src/launcher_config"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var BackendReady bool = false

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	appContext = a.ctx

	cfg, err := launcher_config.InitConfig(appContext) // initialize config
	_ = cfg                                            //ignored
	if err != nil {                                    // if any issue
		BackendReady = false
		runtime.Quit(ctx) // close app
		return
	}
	if err := launcher_accounts.ValidateAll(); err != nil {
		helpers.WarnLauncherLog("[ACCOUNTS] Microsoft account validation failed: " + err.Error())
	}

	BackendReady = true
}

func (a *App) IsBackendReady() bool {
	return BackendReady
}
