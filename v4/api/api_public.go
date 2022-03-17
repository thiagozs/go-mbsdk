package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/thiagozs/go-mbsdk/v4/config"
	"github.com/thiagozs/go-mbsdk/v4/models"
	"github.com/thiagozs/go-mbsdk/v4/pkg/caller"
	"github.com/thiagozs/go-mbsdk/v4/pkg/replacer"
)

func (a *Api) Tickers(symbol string) (models.TickersResponse, error) {
	tickers := models.TickersResponse{}

	c, err := caller.ClientPublic(http.MethodGet, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientPublic")
		}
		return tickers, err
	}

	v, _ := query.Values(models.TickersQuery{Symbols: symbol})
	endpoint, err := replacer.Endpoint(
		replacer.OptKey("TICKERS"),
		replacer.OptSymbol(symbol),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return tickers, err
	}

	bts, err := c.Get(fmt.Sprintf("%s?%s", endpoint, v.Encode()))
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Get")
		}
		return tickers, err
	}

	if err := json.Unmarshal(bts, &tickers); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal Tickers")
		}
		return tickers, err
	}

	return tickers, nil
}
