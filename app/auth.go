package app

import (
	"log"
	"syaniclauncher/src/helpers"
	"syaniclauncher/src/ms_auth"
)

// todo: make a device code version too
func (a *App) AddMicrosoftAccount() {
	authURL, result, err := ms_auth.StartMicrosoftBrowserAuth()
	if err != nil {
		return
	}

	helpers.OpenInBrowser(a.ctx, authURL)

	go func() {
		authResult := <-result
		//todo: implement accounts manager, and do the xbox exchange shi- and store the minecraft account

		if authResult.Err != nil {
			log.Println(authResult.Err)
			return
		}

		log.Println("Microsoft OAuth Successful")
	}()
}
