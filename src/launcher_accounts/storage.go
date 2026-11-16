package launcher_accounts

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"syaniclauncher/src/helpers"

	"gopkg.in/yaml.v3"
)

const accountsNote = "# DO NOT SHARE THIS FILE. IT HOLDS LOGIN CREDENTIALS THAT CAN ACCESS YOUR MINECRAFT ACCOUNT.\n"

type accountFile struct {
	Selected *string   `yaml:"selected"`
	Accounts []Account `yaml:"accounts"`
}

func accountsPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, ".syaniclauncher", "accounts.yml"), nil
}

func load() (accountFile, error) {
	path, err := accountsPath()
	if err != nil {
		return accountFile{}, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return accountFile{Accounts: []Account{}}, nil
	}
	if err != nil {
		return accountFile{}, err
	}

	var accounts accountFile
	if err := yaml.Unmarshal(data, &accounts); err != nil {
		backupPath := fmt.Sprintf("%s.%d.corrupt", path, time.Now().UnixNano())
		if err := os.Rename(path, backupPath); err != nil {
			helpers.WarnLauncherLog(fmt.Sprintf("[ACCOUNTS] corrupt file could not be kept aside: %v", err))
		} else {
			helpers.WarnLauncherLog(fmt.Sprintf("[ACCOUNTS] corrupt file moved to %s", backupPath))
		}
		return accountFile{Accounts: []Account{}}, nil
	}
	if accounts.Accounts == nil {
		accounts.Accounts = []Account{}
	}
	return accounts, nil
}

func save(accounts accountFile) error {
	path, err := accountsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(accounts)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, append([]byte(accountsNote), data...), 0600); err != nil {
		return err
	}
	return os.Chmod(path, 0600)
}

func find(accounts *accountFile, accountUUID string) *Account {
	for i := range accounts.Accounts {
		if accounts.Accounts[i].UUID == accountUUID {
			return &accounts.Accounts[i]
		}
	}
	return nil
}

func upsert(accounts *accountFile, account Account) {
	if existing := find(accounts, account.UUID); existing != nil {
		*existing = account
	} else {
		accounts.Accounts = append(accounts.Accounts, account)
	}
	selected := account.UUID
	accounts.Selected = &selected
}
