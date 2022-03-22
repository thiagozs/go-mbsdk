package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/thiagozs/go-mbsdk/v4/config"
	"github.com/thiagozs/go-mbsdk/v4/models"
	"github.com/thiagozs/go-mbsdk/v4/pkg/caller"
	"github.com/thiagozs/go-mbsdk/v4/pkg/replacer"
)

func (a *Api) GetBalances() (models.ListBalancesResponse, error) {
	balances := models.ListBalancesResponse{}
	errApi := models.ErrorApiResponse{}

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

	res, err := c.GetWithResponse(endpoint)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Get")
		}
		return balances, err
	}
	defer res.Body.Close()

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ReadAll")
		}
		return balances, err
	}

	if config.Config.Debug {
		a.log.Info().
			Str("endpoint", endpoint).
			Int("status_code", res.StatusCode).
			Str("body", string(bts)).
			Msg("")
	}

	if res.StatusCode >= 400 {
		if err := json.Unmarshal(bts, &errApi); err != nil {
			if config.Config.Debug {
				a.log.Error().Stack().Err(err).Msg("Json Unmarshal errApi")
			}
			return balances, err
		}
		return balances, fmt.Errorf("%s - %s", errApi.Code, errApi.Message)
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
	errApi := models.ErrorApiResponse{}

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

	res, err := c.GetWithResponse(endpoint)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Get")
		}
		return acc, err
	}
	defer res.Body.Close()

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ReadAll")
		}
		return acc, err
	}

	if config.Config.Debug {
		a.log.Info().
			Str("endpoint", endpoint).
			Int("status_code", res.StatusCode).
			Str("body", string(bts)).
			Msg("")
	}

	if res.StatusCode >= 400 {
		if err := json.Unmarshal(bts, &errApi); err != nil {
			if config.Config.Debug {
				a.log.Error().Stack().Err(err).Msg("Json Unmarshal errApi")
			}
			return acc, err
		}
		return acc, fmt.Errorf("%s - %s", errApi.Code, errApi.Message)
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
