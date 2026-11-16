package launcher_accounts

import (
	"errors"
	"fmt"
	"time"

	"syaniclauncher/src/ms_auth"
)

func refresh(account *Account) error {
	auth, err := ms_auth.RefreshMicrosoftAuthResult(account.RefreshToken)
	if err != nil {
		return err
	}
	xbox, err := ms_auth.AuthenticateXboxLive(auth.AccessToken)
	if err != nil {
		return err
	}
	xsts, err := ms_auth.AuthorizeXSTS(xbox.Token)
	if err != nil {
		return err
	}
	minecraft, err := ms_auth.AuthenticateMinecraft(xbox.UserHash, xsts.Token)
	if err != nil {
		return err
	}
	profile, err := ms_auth.GetMinecraftProfile(minecraft.AccessToken)
	if err != nil {
		return err
	}

	account.Username = profile.Name
	account.UUID = profile.ID
	account.AccessToken = minecraft.AccessToken
	account.RefreshToken = auth.RefreshToken
	account.ExpiresAt = time.Now().Add(90 * 24 * time.Hour).Unix()
	return nil
}

func GetLaunchAccount(accountUUID string) (*Account, error) {
	accounts, err := load()
	if err != nil {
		return nil, err
	}
	account := find(&accounts, accountUUID)
	if account == nil && accounts.Selected != nil {
		account = find(&accounts, *accounts.Selected)
	}
	if account == nil {
		return nil, nil
	}
	if account.Type == "msa" && account.RefreshToken != "" {
		if err := refresh(account); err != nil {
			return nil, err
		}
		if err := save(accounts); err != nil {
			return nil, err
		}
	}
	result := *account
	return &result, nil
}

func RefreshAll() (AccountList, error) {
	accounts, err := load()
	if err != nil {
		return AccountList{}, err
	}
	var refreshErrors []error
	changed := false
	for i := range accounts.Accounts {
		account := &accounts.Accounts[i]
		if account.Type != "msa" || account.RefreshToken == "" {
			continue
		}
		if err := refresh(account); err != nil {
			refreshErrors = append(refreshErrors, fmt.Errorf("refresh account %q: %w", account.Username, err))
			continue
		}
		changed = true
	}
	if changed {
		if err := save(accounts); err != nil {
			return AccountList{}, err
		}
	}
	return publicList(accounts), errors.Join(refreshErrors...)
}
