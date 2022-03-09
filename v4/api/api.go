package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/shopspring/decimal"
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

func (a *Api) PlaceOrder(kind Kind, symbol, priceIn, pricestopIn, qty string) models.CustomPlaceOrderInfo {
	orderInfo := models.CustomPlaceOrderInfo{}
	order := models.PlaceOrderPayload{Async: true, Type: "limit"}
	price := decimal.RequireFromString(priceIn)
	pricestop := decimal.RequireFromString(pricestopIn)

	cutPrice := strings.Split(price.String(), ".")
	limitPrice, _ := strconv.ParseInt(cutPrice[0], 10, 64)

	cutPriceStop, _ := strconv.ParseInt(pricestop.String(), 10, 64)

	switch kind {
	case BUY:
		order.Side = BUY.String()
	case SELL:
		order.Side = SELL.String()
	case STOP_BUY:
		order.Side = STOP_BUY.String()
		order.StopPrice = int(cutPriceStop)
	case STOP_SELL:
		order.Side = STOP_SELL.String()
		order.StopPrice = int(cutPriceStop)
	}

	order.LimitPrice = int(limitPrice)
	order.Qty = qty

	c, err := caller.ClientWithToken(http.MethodPost, a.cache)
	if err != nil {
		orderInfo.Error = err
		return orderInfo
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("ORDER_PLACE"),
		replacer.OptPair(symbol),
		replacer.OptCache(a.cache))
	if err != nil {
		orderInfo.Error = err
		return orderInfo
	}

	orderInfo.EndPoint = endpoint
	orderInfo.Payload = string(order.ToBytes())

	resp, err := c.PostWithResponse(endpoint, order.ToBytes())
	if err != nil {
		orderInfo.Error = err
		return orderInfo
	}
	defer resp.Body.Close()

	orderInfo.StatusCode = resp.StatusCode

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		orderInfo.Error = err
		return orderInfo
	}

	if resp.StatusCode >= 400 {
		respOrder := models.ErrorPlaceOrderResponse{}
		if err := json.Unmarshal(bts, &respOrder); err != nil {
			orderInfo.Error = err
			return orderInfo
		}
		orderInfo.Error = fmt.Errorf("%s", respOrder.Message)
		orderInfo.StatusCode = resp.StatusCode
		orderInfo.Response = string(respOrder.ToBytes())
		return orderInfo
	}

	respOrder := models.PlaceOrderResponse{}
	if err := json.Unmarshal(bts, &respOrder); err != nil {
		orderInfo.Error = err
		return orderInfo
	}

	if err := a.SetOrder(cutPrice[0], respOrder.OrderID); err != nil {
		orderInfo.Error = err
		return orderInfo
	}

	orderInfo.Response = string(respOrder.ToBytes())
	orderInfo.OrderID = respOrder.OrderID

	return orderInfo
}

func (a *Api) CancelOrder(symbol string, price string) error {

	cutPrice := strings.Split(price, ".")

	c, err := caller.ClientWithToken(http.MethodDelete, a.cache)
	if err != nil {
		return err
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("ORDER_CANCEL"),
		replacer.OptPair(symbol),
		replacer.OptPriceIn(cutPrice[0]),
		replacer.OptCache(a.cache))
	if err != nil {
		return err
	}

	_, err = c.Delete(endpoint, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *Api) CancelAllOrders(symbol string, prices []string) error {

	for _, v := range prices {
		if err := a.CancelOrder(symbol, v); err != nil {
			return err
		}
	}

	return nil
}
