package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/thiagozs/go-mbsdk/v4/config"
	"github.com/thiagozs/go-mbsdk/v4/models"
	"github.com/thiagozs/go-mbsdk/v4/pkg/caller"
	"github.com/thiagozs/go-mbsdk/v4/pkg/replacer"
)

type CandlesOptions func(c *CandlesParameters) error

type CandlesParameters struct {
	Symbols    string `url:"symbol,omitempty"`
	Resolution string `url:"resolution,omitempty"`
	From       int    `url:"from,omitempty"`
	To         int    `url:"to,omitempty"`
	CountBack  int    `url:"countback,omitempty"`
}

func CandSymbols(symbols string) CandlesOptions {
	return func(c *CandlesParameters) error {
		c.Symbols = symbols
		return nil
	}
}

func CandResolution(resolution string) CandlesOptions {
	return func(c *CandlesParameters) error {
		c.Resolution = resolution
		return nil
	}
}

func CandTo(to int) CandlesOptions {
	return func(c *CandlesParameters) error {
		c.To = to
		return nil
	}
}

func CandFrom(from int) CandlesOptions {
	return func(c *CandlesParameters) error {
		c.From = from
		return nil
	}
}

func CandCountBack(countback int) CandlesOptions {
	return func(c *CandlesParameters) error {
		c.CountBack = countback
		return nil
	}
}

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

func (a *Api) OrderBook(symbol, limit string) (models.OrderBookResponse, error) {
	orderbook := models.OrderBookResponse{}
	errApi := models.ErrorApiResponse{}

	c, err := caller.ClientPublic(http.MethodGet, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientPublic")
		}
		return orderbook, err
	}

	endpoint, err := replacer.Endpoint(
		replacer.OptKey("ORDERBOOK"),
		replacer.OptSymbol(symbol),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return orderbook, err
	}

	if len(limit) > 0 {
		v, _ := query.Values(models.OrderBookQuery{Limit: limit})
		endpoint = fmt.Sprintf("%s?%s", endpoint, v.Encode())
	}

	res, err := c.GetWithResponse(endpoint)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Get")
		}
		return orderbook, err
	}
	defer res.Body.Close()

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ReadAll")
		}
		return orderbook, err
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
			return orderbook, err
		}
		return orderbook, fmt.Errorf("%s", errApi.Message)
	}

	if err := json.Unmarshal(bts, &orderbook); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal orderbook")
		}
		return orderbook, err
	}

	return orderbook, nil
}

func (a *Api) Trades(symbol string) (models.TradesResponse, error) {
	trades := models.TradesResponse{}
	errApi := models.ErrorApiResponse{}

	c, err := caller.ClientPublic(http.MethodGet, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientPublic")
		}
		return trades, err
	}

	endpoint, err := replacer.Endpoint(
		replacer.OptKey("TRADES"),
		replacer.OptSymbol(symbol),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return trades, err
	}

	res, err := c.GetWithResponse(endpoint)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Get")
		}
		return trades, err
	}
	defer res.Body.Close()

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ReadAll")
		}
		return trades, err
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
			return trades, err
		}
		return trades, fmt.Errorf("%s", errApi.Message)
	}

	if err := json.Unmarshal(bts, &trades); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal trades")
		}
		return trades, err
	}

	return trades, nil
}

func (a *Api) Symbols(symbol []string) (models.SymbolsResponse, error) {
	symbols := models.SymbolsResponse{}
	errApi := models.ErrorApiResponse{}

	c, err := caller.ClientPublic(http.MethodGet, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientPublic")
		}
		return symbols, err
	}

	endpoint, err := replacer.Endpoint(
		replacer.OptKey("SYMBOLS"),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return symbols, err
	}

	if len(symbol) > 0 {
		v, _ := query.Values(models.SymbolsQuery{Symbols: symbol})
		endpoint = fmt.Sprintf("%s?%s", endpoint, v.Encode())
	}

	res, err := c.GetWithResponse(endpoint)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Get")
		}
		return symbols, err
	}
	defer res.Body.Close()

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ReadAll")
		}
		return symbols, err
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
			return symbols, err
		}
		return symbols, fmt.Errorf("%s", errApi.Message)
	}

	if err := json.Unmarshal(bts, &symbols); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal symbols")
		}
		return symbols, err
	}

	return symbols, nil
}

func (a *Api) Candles(opts ...CandlesOptions) (models.CandlesResponse, error) {
	candles := models.CandlesResponse{}
	errApi := models.ErrorApiResponse{}
	params := &CandlesParameters{}

	for _, op := range opts {
		err := op(params)
		if err != nil {
			return candles, err
		}
	}

	if params.Symbols == "" {
		return candles, fmt.Errorf("parameters 'symbols' is required")
	}

	if params.Resolution == "" {
		return candles, fmt.Errorf("parameters 'resolution' is required")
	}

	if params.To <= 0 {
		return candles, fmt.Errorf("parameters 'to' is required")
	}

	if params.From <= 0 {
		return candles, fmt.Errorf("parameters 'from' is required")
	}

	c, err := caller.ClientPublic(http.MethodGet, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientPublic")
		}
		return candles, err
	}

	endpoint, err := replacer.Endpoint(
		replacer.OptKey("CANDLES"),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return candles, err
	}

	v, _ := query.Values(params)
	endpoint = fmt.Sprintf("%s?%s", endpoint, v.Encode())

	res, err := c.GetWithResponse(endpoint)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Get")
		}
		return candles, err
	}
	defer res.Body.Close()

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ReadAll")
		}
		return candles, err
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
			return candles, err
		}
		return candles, fmt.Errorf("%s", errApi.Message)
	}

	if err := json.Unmarshal(bts, &candles); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal candles")
		}
		return candles, err
	}

	return candles, nil
}
