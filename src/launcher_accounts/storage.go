package launcher_accounts

import (
	"fmt"
	"os"
	"path/filepath"

	"syaniclauncher/src/helpers"
	"syaniclauncher/src/helpers_yml"
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

	var accounts accountFile
	exists, err := helpers_yml.Read(path, &accounts)
	if !exists {
		return accountFile{Accounts: []Account{}}, nil
	}
	if err != nil {
		if backupPath, quarantineErr := helpers_yml.Quarantine(path); quarantineErr != nil {
			helpers.WarnLauncherLog(fmt.Sprintf("[ACCOUNTS] corrupt file could not be kept aside: %v", quarantineErr))
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
	data, err := helpers_yml.Marshal(accounts)
	if err != nil {
		return err
	}
	return helpers_yml.Write(path, append([]byte(accountsNote), data...), 0600)
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
