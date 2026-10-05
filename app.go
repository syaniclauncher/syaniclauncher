package main

import (
	"context"
	"fmt"
	"syaniclauncher/src/config"
	"syaniclauncher/src/helpers"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Wails default stuff
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	helpers.SetContext(a.ctx) // make ctx global

	cfg, err := config.InitConfig() // initialize config
	_ = cfg                         //ignored
	if err != nil {                 // if any issue
		config.BackendReady = false
		runtime.Quit(ctx) // close app
		return
	}

	config.BackendReady = true
}

func (a *App) IsBackendReady() bool { // first letter must be cap btw to register
	return config.BackendReady
}

// test: CheckConnection function to check connection between backend and frontend
func (a *App) CheckConnection() string {
	fmt.Println("Connection Check Recived")
	return fmt.Sprintln("Absolute good :)")
}
