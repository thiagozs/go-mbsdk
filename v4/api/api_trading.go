package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/thiagozs/go-mbsdk/v4/models"
	"github.com/thiagozs/go-mbsdk/v4/pkg/caller"
	"github.com/thiagozs/go-mbsdk/v4/pkg/replacer"

	"github.com/google/go-querystring/query"
)

type OrdersParams func(o *OrdersPameters) error

type OrdersPameters struct {
	has_executions  string `url:"has_executions,omitempty"`
	side            string `url:"side,omitempty"`
	status          string `url:"status,omitempty"`
	id_from         string `url:"id_from,omitempty"`
	id_to           string `url:"id_to,omitempty"`
	created_at_from string `url:"created_at_from,omitempty"`
	created_at_to   string `url:"created_at_to,omitempty"`
}

func OdrHasExec(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.has_executions = value
		return nil
	}
}

func OrdSide(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.side = value
		return nil
	}
}

func OrdSatus(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.status = value
		return nil
	}
}

func OrdIdFrom(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.id_from = value
		return nil
	}
}

func OrdIdTo(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.id_to = value
		return nil
	}
}

func OrdCreatedFrom(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.created_at_from = value
		return nil
	}
}

func OrdCreatedTo(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.created_at_to = value
		return nil
	}
}

func (a *Api) PlaceOrder(kind Kind, symbol, priceIn, pricestopIn, qty string) models.CustomPlaceOrderInfo {
	orderInfo := models.CustomPlaceOrderInfo{}
	order := models.PlaceOrderPayload{Async: true, Type: "limit"}

	price, _ := decimal.NewFromString(priceIn)

	pricestop, _ := decimal.NewFromString(pricestopIn)

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

	if price.GreaterThan(decimal.RequireFromString("0")) {
		order.LimitPrice = int(limitPrice)
	}

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

	if err := a.SetOrder(models.OrdersIndex{
		Symbol: symbol,
		ID:     respOrder.OrderID,
		Price:  price.String(),
		Side:   order.Side,
		Type:   order.Type,
	}); err != nil {
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

func (a *Api) CancelAllOpenOrders(symbol string) error {

	c, err := caller.ClientWithToken(http.MethodDelete, a.cache)
	if err != nil {
		return err
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("ORDER_CANCEL_ALL"),
		replacer.OptPair(symbol))

	if err != nil {
		return err
	}

	_, err = c.Delete(endpoint, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *Api) GetOrder(symbol string) (models.GetOrderResponse, error) {
	order := models.GetOrderResponse{}
	c, err := caller.ClientWithToken(http.MethodGet, a.cache)
	if err != nil {
		return order, err
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("ORDER_GET"),
		replacer.OptPair(symbol),
		replacer.OptCache(a.cache),
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

func (a *Api) ListOrders(symbol string, opts ...OrdersParams) (models.ListOrderResponse, error) {
	order := models.ListOrderResponse{}
	params := &OrdersPameters{}

	for _, op := range opts {
		err := op(params)
		if err != nil {
			return order, err
		}
	}

	c, err := caller.ClientWithToken(http.MethodGet, a.cache)
	if err != nil {
		return order, err
	}

	// gen parametes
	v, _ := query.Values(params)

	endpoint, err := replacer.Endpoint(replacer.OptKey("ORDER_LIST"),
		replacer.OptPair(symbol),
		replacer.OptCache(a.cache),
		replacer.OptParams(v.Encode()),
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
