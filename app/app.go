package app

import (
	"context"
	"sync"
)

var appContext context.Context
var BackendReady bool = false

// Wails default stuff
type App struct {
	ctx           context.Context
	authMu        sync.Mutex
	microsoftAuth *microsoftAuthSession
}

func NewApp() *App {
	return &App{}
}
