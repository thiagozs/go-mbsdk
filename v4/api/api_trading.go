package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/thiagozs/go-mbsdk/v4/config"
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
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientWithToken")
		}
		orderInfo.Error = err
		return orderInfo
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("ORDER_PLACE"),
		replacer.OptSymbol(symbol),
		replacer.OptCache(a.cache),
		replacer.OptLog(a.log),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		orderInfo.Error = err
		return orderInfo
	}

	orderInfo.EndPoint = endpoint
	orderInfo.Payload = string(order.ToBytes())

	resp, err := c.PostWithResponse(endpoint, order.ToBytes())
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("PostWithResponse")
		}
		orderInfo.Error = err
		return orderInfo
	}
	defer resp.Body.Close()

	orderInfo.StatusCode = resp.StatusCode

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ReadAll")
		}
		orderInfo.Error = err
		return orderInfo
	}

	if resp.StatusCode >= 400 {
		respOrder := models.ErrorPlaceOrderResponse{}
		if err := json.Unmarshal(bts, &respOrder); err != nil {
			if config.Config.Debug {
				a.log.Error().Stack().Err(err).Msg("Json Unmarshal ErrorPlaceOrderResponse")
			}
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
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal PlaceOrderResponse")
		}
		orderInfo.Error = err
		return orderInfo
	}

	if err := a.SetOrder([]models.OrdersIndex{
		{
			Symbol: symbol,
			ID:     respOrder.OrderID,
			Price:  price.String(),
			Side:   order.Side,
			Type:   order.Type,
		},
	}); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Cache SetOrder")
		}
		orderInfo.Error = err
		return orderInfo
	}

	orderInfo.Response = string(respOrder.ToBytes())
	orderInfo.OrderID = respOrder.OrderID

	return orderInfo
}

func (a *Api) CancelOrder(symbol string, id string) error {

	c, err := caller.ClientWithToken(http.MethodDelete, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientWithToken")
		}
		return err
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("ORDER_CANCEL"),
		replacer.OptSymbol(symbol),
		replacer.OptCache(a.cache),
		replacer.OptOrderId(id),
		replacer.OptLog(a.log),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return err
	}

	_, err = c.Delete(endpoint, nil)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Delete")
		}
		return err
	}

	return nil
}

func (a *Api) CancelAllOrders(symbol string) error {

	val, err := a.cache.GetKeyVal(config.ORDERS_INDEX.String())
	if err != nil {
		return err
	}

	ordersIndex := []models.OrdersIndex{}
	if err := json.Unmarshal([]byte(val), &ordersIndex); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal OrderIndex")
		}
		return err
	}

	for i, v := range ordersIndex {
		if strings.EqualFold(v.Symbol, symbol) {
			if err := a.CancelOrder(symbol, v.ID); err != nil {
				if config.Config.Debug {
					a.log.Error().Stack().Err(err).Msg("CancelOrder")
				}
				return err
			}
			ordersIndex = append(ordersIndex[:i], ordersIndex[i+1:]...)
		}
	}

	if err := a.SetOrder(ordersIndex); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("SetOrder")
		}
		return err
	}

	return nil
}

func (a *Api) CancelAllOpenOrders(symbol string) error {

	c, err := caller.ClientWithToken(http.MethodDelete, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientWithToken")
		}
		return err
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("ORDER_CANCEL_ALL"),
		replacer.OptSymbol(symbol),
		replacer.OptLog(a.log),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return err
	}

	_, err = c.Delete(endpoint, nil)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Delete")
		}
		return err
	}

	return nil
}

func (a *Api) GetOrder(symbol string) (models.GetOrderResponse, error) {
	order := models.GetOrderResponse{}
	c, err := caller.ClientWithToken(http.MethodGet, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientWithToken")
		}
		return order, err
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("ORDER_GET"),
		replacer.OptSymbol(symbol),
		replacer.OptCache(a.cache),
		replacer.OptLog(a.log),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return order, err
	}

	bts, err := c.Get(endpoint)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Get")
		}
		return order, err
	}

	if err := json.Unmarshal(bts, &order); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal GetOrderResponse")
		}
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
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientWithToken")
		}
		return order, err
	}

	// gen parametes
	v, _ := query.Values(params)
	endpoint, err := replacer.Endpoint(replacer.OptKey("ORDER_LIST"),
		replacer.OptSymbol(symbol),
		replacer.OptCache(a.cache),
		replacer.OptLog(a.log),
		replacer.OptParams(v.Encode()),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return order, err
	}

	bts, err := c.Get(endpoint)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Get")
		}
		return order, err
	}

	if err := json.Unmarshal(bts, &order); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal ListOrderResponse")
		}
		return order, err
	}

	return order, nil
}
