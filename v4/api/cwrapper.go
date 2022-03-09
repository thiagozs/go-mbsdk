package api

import (
	"github.com/thiagozs/go-mbsdk/v4/config"
	"github.com/thiagozs/go-mbsdk/v4/models"
)

func (a *Api) SetBalance(balance models.GetBalancesResponse) error {
	return a.cache.SetKeyValAsJSON(config.BALANCE.String(), balance)
}

func (a *Api) SetAccounts(acc models.GetAccountsResponse) error {
	return a.cache.SetKeyValAsJSON(config.ACCOUNTS.String(), acc)
}

func (a *Api) SetAuthorize(auth models.AuthoritionToken) error {
	return a.cache.SetKeyValAsJSON(config.AUTHORIZE.String(), auth)
}

func (a *Api) SetOrder(key, value string) error {
	return a.cache.SetKeyVal(key, value)
}

func (a *Api) GetOrder(key string) (string, error) {
	return a.cache.GetKeyVal(key)
}

func (a *Api) DeleteOrder(key string) (string, error) {
	return a.cache.DeleteKey(key)
}

func (a *Api) GetKeyVal(key string) (string, error) {
	return a.cache.GetKeyVal(key)
}
