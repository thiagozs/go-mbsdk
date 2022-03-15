package api

import (
	"encoding/json"
	"net/http"

	"github.com/thiagozs/go-mbsdk/v4/models"
	"github.com/thiagozs/go-mbsdk/v4/pkg/caller"
	"github.com/thiagozs/go-mbsdk/v4/pkg/replacer"
)

func (a *Api) GetBalances() (models.ListBalancesResponse, error) {
	balances := models.ListBalancesResponse{}
	c, err := caller.ClientWithToken(http.MethodGet, a.cache)
	if err != nil {
		return balances, err
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("BALANCE_LIST"),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		return balances, err
	}

	bts, err := c.Get(endpoint)
	if err != nil {
		return balances, err
	}

	if err := json.Unmarshal(bts, &balances); err != nil {
		return balances, err
	}

	if err := a.SetBalance(balances); err != nil {
		return balances, err
	}

	return balances, nil
}

func (a *Api) GetAccounts() (models.ListAccountsResponse, error) {

	acc := models.ListAccountsResponse{}

	c, err := caller.ClientWithToken(http.MethodGet, a.cache)
	if err != nil {
		return acc, err
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("ACCOUNTS"),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		return acc, err
	}

	bts, err := c.Get(endpoint)
	if err != nil {
		return acc, err
	}

	if err := json.Unmarshal(bts, &acc); err != nil {
		return acc, err
	}

	if err := a.SetAccounts(acc); err != nil {
		return acc, err
	}

	return acc, nil
}
