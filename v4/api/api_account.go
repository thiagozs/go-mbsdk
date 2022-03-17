package api

import (
	"encoding/json"
	"net/http"

	"github.com/thiagozs/go-mbsdk/v4/config"
	"github.com/thiagozs/go-mbsdk/v4/models"
	"github.com/thiagozs/go-mbsdk/v4/pkg/caller"
	"github.com/thiagozs/go-mbsdk/v4/pkg/replacer"
)

func (a *Api) GetBalances() (models.ListBalancesResponse, error) {
	balances := models.ListBalancesResponse{}
	c, err := caller.ClientWithToken(http.MethodGet, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientWithToken")
		}
		return balances, err
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("BALANCE_LIST"),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return balances, err
	}

	bts, err := c.Get(endpoint)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Get")
		}
		return balances, err
	}

	if err := json.Unmarshal(bts, &balances); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal Balances")
		}
		return balances, err
	}

	if err := a.CacheSetBalance(balances); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("CacheSetBalance")
		}
		return balances, err
	}

	return balances, nil
}

func (a *Api) GetAccounts() (models.ListAccountsResponse, error) {

	acc := models.ListAccountsResponse{}

	c, err := caller.ClientWithToken(http.MethodGet, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientWithToken")
		}
		return acc, err
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("ACCOUNTS"),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return acc, err
	}

	bts, err := c.Get(endpoint)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Get")
		}
		return acc, err
	}

	if err := json.Unmarshal(bts, &acc); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal ListAccountsResponse")
		}
		return acc, err
	}

	if err := a.CacheSetAccounts(acc); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("SetAccounts")
		}
		return acc, err
	}

	return acc, nil
}
