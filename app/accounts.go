package app

import "syaniclauncher/src/launcher_accounts"

func (a *App) AccountsList() (launcher_accounts.AccountList, error) {
	return launcher_accounts.ListAccounts()
}

func (a *App) AccountsRefresh() (launcher_accounts.AccountList, error) {
	return launcher_accounts.RefreshAll()
}

func (a *App) AccountAddOffline(username string) (launcher_accounts.AccountList, error) {
	return launcher_accounts.AddOffline(username)
}

func (a *App) AccountRemove(accountUUID string) (launcher_accounts.AccountList, error) {
	return launcher_accounts.Remove(accountUUID)
}

func (a *App) AccountSetSelected(accountUUID string) (launcher_accounts.AccountList, error) {
	return launcher_accounts.SetSelected(accountUUID)
}
