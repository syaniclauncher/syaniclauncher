package launcher_accounts

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"syaniclauncher/src/helpers"
	"syaniclauncher/src/launcher_cache"
)

type Account struct {
	Type         string `yaml:"type" json:"type"`
	Username     string `yaml:"username" json:"username"`
	UUID         string `yaml:"uuid" json:"uuid"`
	AccessToken  string `yaml:"access_token,omitempty" json:"-"`
	RefreshToken string `yaml:"refresh_token,omitempty" json:"-"`
	ExpiresAt    int64  `yaml:"expires_at,omitempty" json:"-"`
	Valid        bool   `yaml:"valid" json:"-"`
}

type PublicAccount struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Expired  *bool  `json:"expired,omitempty"`
	Head     string `json:"head,omitempty"`
	Valid    bool   `json:"valid"`
}

type AccountList struct {
	Selected *string         `json:"selected"`
	Accounts []PublicAccount `json:"accounts"`
}

func publicList(accounts accountFile) AccountList {
	result := AccountList{
		Selected: accounts.Selected,
		Accounts: make([]PublicAccount, 0, len(accounts.Accounts)),
	}
	now := time.Now().Unix()
	for _, account := range accounts.Accounts {
		public := PublicAccount{
			Type:     account.Type,
			Username: account.Username,
			UUID:     account.UUID,
			Valid:    account.Valid,
		}
		if account.Type == "msa" {
			expired := !account.Valid || account.ExpiresAt <= now
			public.Expired = &expired
			if account.Valid {
				head, err := launcher_cache.GetPlayerHead(account.UUID)
				if err == nil {
					public.Head = "data:image/png;base64," + base64.StdEncoding.EncodeToString(head)
				}
			}
		}
		result.Accounts = append(result.Accounts, public)
	}
	return result
}

func ListAccounts() (AccountList, error) {
	accounts, err := load()
	if err != nil {
		return AccountList{}, err
	}
	return publicList(accounts), nil
}

func AddOffline(username string, skipChecks bool) (AccountList, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		username = "Player"
	}
	if !skipChecks && !helpers.IsValidMinecraftUsername(username) {
		return AccountList{}, fmt.Errorf("Username must be 3-16 characters and use only letters, numbers, and underscores")
	}
	accounts, err := load()
	if err != nil {
		return AccountList{}, err
	}
	upsert(&accounts, Account{
		Type:     "offline",
		Username: username,
		UUID:     helpers.OfflineUUID(username),
		Valid:    true,
	})
	if err := save(accounts); err != nil {
		return AccountList{}, err
	}
	return publicList(accounts), nil
}

func AddMSA(account Account) (AccountList, error) {
	if strings.TrimSpace(account.Username) == "" || strings.TrimSpace(account.UUID) == "" {
		return AccountList{}, fmt.Errorf("Microsoft account username and UUID are required")
	}
	account.Type = "msa"
	account.Valid = true
	accounts, err := load()
	if err != nil {
		return AccountList{}, err
	}
	upsert(&accounts, account)
	if err := save(accounts); err != nil {
		return AccountList{}, err
	}
	return publicList(accounts), nil
}

func Remove(accountUUID string) (AccountList, error) {
	accounts, err := load()
	if err != nil {
		return AccountList{}, err
	}
	remaining := accounts.Accounts[:0]
	for _, account := range accounts.Accounts {
		if account.UUID != accountUUID {
			remaining = append(remaining, account)
		}
	}
	accounts.Accounts = remaining
	if accounts.Selected != nil && *accounts.Selected == accountUUID {
		accounts.Selected = nil
		if len(accounts.Accounts) > 0 {
			selected := accounts.Accounts[0].UUID
			accounts.Selected = &selected
		}
	}
	if err := save(accounts); err != nil {
		return AccountList{}, err
	}
	return publicList(accounts), nil
}

func SetSelected(accountUUID string) (AccountList, error) {
	accounts, err := load()
	if err != nil {
		return AccountList{}, err
	}
	if find(&accounts, accountUUID) != nil {
		selected := accountUUID
		accounts.Selected = &selected
		if err := save(accounts); err != nil {
			return AccountList{}, err
		}
	}
	return publicList(accounts), nil
}
