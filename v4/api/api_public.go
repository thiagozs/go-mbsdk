package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/thiagozs/go-mbsdk/v4/models"
	"github.com/thiagozs/go-mbsdk/v4/pkg/caller"
	"github.com/thiagozs/go-mbsdk/v4/pkg/replacer"
	"github.com/thiagozs/go-mbsdk/v4/pkg/utils"
)

func (a *Api) Tickers(symbol string) (models.TickersResponse, error) {
	tickers := models.TickersResponse{}

	c, err := caller.ClientPublic(http.MethodGet, a.cache)
	if err != nil {
		return tickers, err
	}

	pair, _ := utils.PairQuote(symbol)

	v, _ := query.Values(models.TickersQuery{Symbols: symbol})
	endpoint, err := replacer.Endpoint(
		replacer.OptKey("TICKERS"),
		replacer.OptPair(pair),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		return tickers, err
	}

	bts, err := c.Get(fmt.Sprintf("%s?%s", endpoint, v.Encode()))
	if err != nil {
		return tickers, err
	}

	if err := json.Unmarshal(bts, &tickers); err != nil {
		return tickers, err
	}

	return tickers, nil
}
