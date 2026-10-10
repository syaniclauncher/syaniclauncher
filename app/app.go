package app

import (
	"context"
)

var appContext context.Context
var BackendReady bool = false

// Wails default stuff
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}
