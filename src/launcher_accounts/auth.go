package launcher_accounts

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"syaniclauncher/src/ms_auth"
)

func isTransientAuthError(err error) bool {
	if err == nil {
		return false
	}
	var networkError net.Error
	if errors.As(err, &networkError) {
		return true
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "http 408") ||
		strings.Contains(message, "http 429") ||
		strings.Contains(message, "http 5")
}

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
	account.Valid = true
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
	if account.Type == "msa" {
		if account.RefreshToken != "" {
			if err := refresh(account); err != nil {
				if !isTransientAuthError(err) {
					account.Valid = false
					_ = save(accounts)
				}
				return nil, err
			}
		} else {
			account.Valid = false
			_ = save(accounts)
			return nil, fmt.Errorf("Microsoft refresh token is empty")
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
		if account.Type != "msa" {
			continue
		}
		if account.RefreshToken == "" {
			account.Valid = false
			changed = true
			refreshErrors = append(refreshErrors, fmt.Errorf("refresh account %q: Microsoft refresh token is empty", account.Username))
			continue
		}
		if err := refresh(account); err != nil {
			if !isTransientAuthError(err) {
				account.Valid = false
				changed = true
			}
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

func ValidateAll() error {
	accounts, err := load()
	if err != nil {
		return err
	}

	var validationErrors []error
	changed := false
	for i := range accounts.Accounts {
		account := &accounts.Accounts[i]
		if account.Type != "msa" {
			continue
		}
		if account.RefreshToken == "" {
			account.Valid = false
			changed = true
			validationErrors = append(validationErrors, fmt.Errorf("account %q: Microsoft refresh token is empty", account.Username))
			continue
		}
		if err := refresh(account); err != nil {
			if !isTransientAuthError(err) {
				account.Valid = false
				changed = true
			}
			validationErrors = append(validationErrors, fmt.Errorf("account %q: %w", account.Username, err))
			continue
		}
		changed = true
	}
	if changed {
		if err := save(accounts); err != nil {
			return err
		}
	}
	return errors.Join(validationErrors...)
}
