package main

import (
	"context"
	"fmt"
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
	// put init here (will be commited soon)
}

func (a *App) IsBackendReady() bool { // first letter must be cap btw to register
	// check if init checks are completed via a global variable ig
	// send false if init checks not done, true if done

	// frontend would loop until backendReady event is received
	// or maybe close the app if error? that would be somewhere else
	return true
}

// test: CheckConnection function to check connection between backend and frontend
func (a *App) CheckConnection() string {
	fmt.Println("Connection Check Recived")
	return fmt.Sprintln("Absolute good :)")
}
