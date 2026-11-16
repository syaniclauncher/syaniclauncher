package app

import (
	"fmt"
	"time"

	"syaniclauncher/src/helpers"
	"syaniclauncher/src/launcher_accounts"
	"syaniclauncher/src/ms_auth"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const accountAuthEvent = "accounts:auth"

type MicrosoftDeviceLogin struct {
	VerificationURI string `json:"verification_uri"`
	UserCode        string `json:"user_code"`
}

func (a *App) StartMicrosoftBrowserLogin() error {
	authURL, result, err := ms_auth.StartMicrosoftBrowserAuth()
	if err != nil {
		return err
	}
	helpers.OpenInBrowser(a.ctx, authURL)
	go a.finishMicrosoftLogin(result)
	return nil
}

func (a *App) StartMicrosoftDeviceLogin() (MicrosoftDeviceLogin, error) {
	result, err, verificationURI, userCode := ms_auth.StartMicrosoftDeviceCodeAuth()
	if err != nil {
		return MicrosoftDeviceLogin{}, err
	}
	go a.finishMicrosoftLogin(result)
	return MicrosoftDeviceLogin{
		VerificationURI: verificationURI,
		UserCode:        userCode,
	}, nil
}

func (a *App) finishMicrosoftLogin(result <-chan ms_auth.AuthResult) {
	auth, ok := <-result
	if !ok {
		a.emitAccountAuthError(fmt.Errorf("Microsoft sign-in ended unexpectedly"))
		return
	}
	if auth.Err != nil {
		a.emitAccountAuthError(auth.Err)
		return
	}

	xbox, err := ms_auth.AuthenticateXboxLive(auth.AccessToken)
	if err != nil {
		a.emitAccountAuthError(err)
		return
	}
	xsts, err := ms_auth.AuthorizeXSTS(xbox.Token)
	if err != nil {
		a.emitAccountAuthError(err)
		return
	}
	minecraft, err := ms_auth.AuthenticateMinecraft(xbox.UserHash, xsts.Token)
	if err != nil {
		a.emitAccountAuthError(err)
		return
	}
	profile, err := ms_auth.GetMinecraftProfile(minecraft.AccessToken)
	if err != nil {
		a.emitAccountAuthError(err)
		return
	}

	_, err = launcher_accounts.AddMSA(launcher_accounts.Account{
		Username:     profile.Name,
		UUID:         profile.ID,
		AccessToken:  minecraft.AccessToken,
		RefreshToken: auth.RefreshToken,
		ExpiresAt:    time.Now().Add(90 * 24 * time.Hour).Unix(),
	})
	if err != nil {
		a.emitAccountAuthError(err)
		return
	}
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, accountAuthEvent, "success", "")
	}
}

func (a *App) emitAccountAuthError(err error) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, accountAuthEvent, "error", err.Error())
	}
}
