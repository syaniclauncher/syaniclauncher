package app

import (
	"context"
	"errors"
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

type microsoftAuthSession struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (a *App) beginMicrosoftAuth() (*microsoftAuthSession, error) {
	if a.ctx == nil {
		return nil, fmt.Errorf("app is not initialized")
	}
	ctx, cancel := context.WithCancel(a.ctx)
	session := &microsoftAuthSession{ctx: ctx, cancel: cancel}

	a.authMu.Lock()
	previous := a.microsoftAuth
	a.microsoftAuth = session
	a.authMu.Unlock()
	if previous != nil {
		previous.cancel()
	}
	return session, nil
}

func (a *App) finishMicrosoftAuth(session *microsoftAuthSession) {
	session.cancel()
	a.authMu.Lock()
	if a.microsoftAuth == session {
		a.microsoftAuth = nil
	}
	a.authMu.Unlock()
}

func (a *App) CancelMicrosoftLogin() {
	a.authMu.Lock()
	session := a.microsoftAuth
	a.microsoftAuth = nil
	a.authMu.Unlock()
	if session != nil {
		session.cancel()
	}
}

func (a *App) StartMicrosoftBrowserLogin() (string, error) {
	session, err := a.beginMicrosoftAuth()
	if err != nil {
		return "", err
	}
	authURL, result, err := ms_auth.StartMicrosoftBrowserAuthContext(session.ctx)
	if err != nil {
		a.finishMicrosoftAuth(session)
		return "", err
	}
	if err := session.ctx.Err(); err != nil {
		a.finishMicrosoftAuth(session)
		return "", err
	}
	helpers.OpenInBrowser(a.ctx, authURL)
	go a.finishMicrosoftLogin(session, result)
	return authURL, nil
}

func (a *App) StartMicrosoftDeviceLogin() (MicrosoftDeviceLogin, error) {
	session, err := a.beginMicrosoftAuth()
	if err != nil {
		return MicrosoftDeviceLogin{}, err
	}
	result, err, verificationURI, userCode := ms_auth.StartMicrosoftDeviceCodeAuthContext(session.ctx)
	if err != nil {
		a.finishMicrosoftAuth(session)
		return MicrosoftDeviceLogin{}, err
	}
	if err := session.ctx.Err(); err != nil {
		a.finishMicrosoftAuth(session)
		return MicrosoftDeviceLogin{}, err
	}
	go a.finishMicrosoftLogin(session, result)
	return MicrosoftDeviceLogin{
		VerificationURI: verificationURI,
		UserCode:        userCode,
	}, nil
}

func (a *App) finishMicrosoftLogin(session *microsoftAuthSession, result <-chan ms_auth.AuthResult) {
	defer a.finishMicrosoftAuth(session)

	var auth ms_auth.AuthResult
	var ok bool
	select {
	case <-session.ctx.Done():
		return
	case auth, ok = <-result:
	}
	if session.ctx.Err() != nil {
		return
	}
	if !ok {
		a.emitAccountAuthError(session.ctx, fmt.Errorf("Microsoft sign-in ended unexpectedly"))
		return
	}
	if auth.Err != nil {
		a.emitAccountAuthError(session.ctx, auth.Err)
		return
	}

	xbox, err := ms_auth.AuthenticateXboxLiveContext(session.ctx, auth.AccessToken)
	if err != nil {
		a.emitAccountAuthError(session.ctx, err)
		return
	}
	xsts, err := ms_auth.AuthorizeXSTSContext(session.ctx, xbox.Token)
	if err != nil {
		a.emitAccountAuthError(session.ctx, err)
		return
	}
	minecraft, err := ms_auth.AuthenticateMinecraftContext(session.ctx, xbox.UserHash, xsts.Token)
	if err != nil {
		a.emitAccountAuthError(session.ctx, err)
		return
	}
	profile, err := ms_auth.GetMinecraftProfileContext(session.ctx, minecraft.AccessToken)
	if err != nil {
		a.emitAccountAuthError(session.ctx, err)
		return
	}
	if session.ctx.Err() != nil {
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
		a.emitAccountAuthError(session.ctx, err)
		return
	}
	if a.ctx != nil && session.ctx.Err() == nil {
		runtime.EventsEmit(a.ctx, accountAuthEvent, "success", "")
	}
}

func (a *App) emitAccountAuthError(ctx context.Context, err error) {
	if a.ctx != nil && ctx.Err() == nil && !errors.Is(err, context.Canceled) {
		runtime.EventsEmit(a.ctx, accountAuthEvent, "error", err.Error())
	}
}
