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
	"github.com/thiagozs/go-mbsdk/v4/pkg/utils"
)

func New(opts ...Options) (*Api, error) {

	mts := &ApiCfg{}

	for _, op := range opts {
		err := op(mts)
		if err != nil {
			return &Api{}, err
		}
	}
	config.Config.Login = mts.key
	config.Config.Password = mts.secret
	config.Config.Cache = mts.cache

	return &Api{mts.cache}, nil
}

func (a *Api) AuthorizationToken() (models.AuthoritionToken, error) {
	auth := models.AuthoritionToken{}
	c, err := caller.ClientWithForm(http.MethodPost, a.cache)
	if err != nil {
		return auth, err
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("AUTHORIZE"),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		return auth, err
	}

	res, err := c.PostFormWithResponse(endpoint)
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()
	if err != nil {
		return auth, err
	}

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return auth, err
	}

	if err := json.Unmarshal(bts, &auth); err != nil {
		return auth, err
	}

	if err := a.SetAuthorize(auth); err != nil {
		return auth, err
	}

	return auth, nil
}

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

func (a *Api) GetAccounts() (models.GetAccountsResponse, error) {

	acc := models.GetAccountsResponse{}

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

func (a *Api) GetBalances() (models.GetBalancesResponse, error) {
	balances := models.GetBalancesResponse{}
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

	fmt.Printf("result = %s\n", bts)

	if err := json.Unmarshal(bts, &balances); err != nil {
		return balances, err
	}

	if err := a.SetBalance(balances); err != nil {
		return balances, err
	}

	return balances, nil
}

func (a *Api) GetOrderInfo(symbol, price string) (models.GetOrderResponse, error) {
	order := models.GetOrderResponse{}
	c, err := caller.ClientWithToken(http.MethodGet, a.cache)
	if err != nil {
		return order, err
	}

	pair, _ := utils.PairQuote(symbol)

	endpoint, err := replacer.Endpoint(replacer.OptKey("ORDER_GET"),
		replacer.OptPair(pair),
		replacer.OptCache(a.cache),
		replacer.OptPriceIn(price),
	)
	if err != nil {
		return order, err
	}

	bts, err := c.Get(endpoint)
	if err != nil {
		return order, err
	}

	if err := json.Unmarshal(bts, &order); err != nil {
		return order, err
	}

	return order, nil
}
