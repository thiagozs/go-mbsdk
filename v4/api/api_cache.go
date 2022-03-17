package api

import (
	"encoding/json"

	"github.com/thiagozs/go-mbsdk/v4/config"
	"github.com/thiagozs/go-mbsdk/v4/models"
)

func (a *Api) CacheSetBalance(balance models.ListBalancesResponse) error {
	return a.cache.SetKeyValAsJSON(config.BALANCE.String(), balance)
}

func (a *Api) CacheSetAccounts(acc models.ListAccountsResponse) error {
	return a.cache.SetKeyValAsJSON(config.ACCOUNTS.String(), acc)
}

func (a *Api) CacheSetAuthorize(auth models.AuthoritionToken) error {
	return a.cache.SetKeyValAsJSON(config.AUTHORIZE.String(), auth)
}

func (a *Api) CacheSetOrder(orderIn []models.OrdersIndex) error {
	orders := []models.OrdersIndex{}

	val, err := a.cache.GetKeyVal(config.ORDERS_INDEX.String())
	if err != nil && err.Error() != "not found" {
		return err
	}
	if len(val) > 0 {
		if err := json.Unmarshal([]byte(val), &orders); err != nil {
			return err
		}
	}
	orders = append(orders, orderIn...)
	return a.cache.SetKeyValAsJSON(config.ORDERS_INDEX.String(), orders)
}

func (a *Api) CacheGetOrder() (string, error) {
	return a.cache.GetKeyVal(config.ORDERS_INDEX.String())
}

func (a *Api) CacheDeleteOrder() (string, error) {
	return a.cache.DeleteKey(config.ORDERS_INDEX.String())
}

func (a *Api) GetKeyVal(key string) (string, error) {
	return a.cache.GetKeyVal(key)
}
