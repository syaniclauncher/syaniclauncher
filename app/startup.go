package app

import (
	"context"
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

	BackendReady = true
}

func (a *App) IsBackendReady() bool {
	return BackendReady
}
