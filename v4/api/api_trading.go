package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/thiagozs/go-cache/v1/cache/drivers/kind"
	"github.com/thiagozs/go-mbsdk/v4/config"
	"github.com/thiagozs/go-mbsdk/v4/models"
	"github.com/thiagozs/go-mbsdk/v4/pkg/caller"
	"github.com/thiagozs/go-mbsdk/v4/pkg/replacer"

	"github.com/google/go-querystring/query"
)

type OrdersParams func(o *OrdersPameters) error

type PlaceOrdersParams func(o *PlaceOrdersPameters) error

type PlaceOrdersPameters struct {
	Symbol    string
	Side      string
	Price     string
	Kind      Kind
	Type      string
	PriceStop string
	Quantity  string
}

func PoSymbol(value string) PlaceOrdersParams {
	return func(a *PlaceOrdersPameters) error {
		a.Symbol = value
		return nil
	}
}

func PoSide(value string) PlaceOrdersParams {
	return func(a *PlaceOrdersPameters) error {
		a.Side = value
		return nil
	}
}

func PoPrice(value string) PlaceOrdersParams {
	return func(a *PlaceOrdersPameters) error {
		a.Price = value
		return nil
	}
}

func PoKind(value Kind) PlaceOrdersParams {
	return func(a *PlaceOrdersPameters) error {
		a.Kind = value
		return nil
	}
}

func PoType(value string) PlaceOrdersParams {
	return func(a *PlaceOrdersPameters) error {
		a.Type = value
		return nil
	}
}

func PoPriceStop(value string) PlaceOrdersParams {
	return func(a *PlaceOrdersPameters) error {
		a.PriceStop = value
		return nil
	}
}

func PoQty(value string) PlaceOrdersParams {
	return func(a *PlaceOrdersPameters) error {
		a.Quantity = value
		return nil
	}
}

type OrdersPameters struct {
	HasExecutions string `url:"has_executions,omitempty"`
	Side          string `url:"side,omitempty"`
	Status        string `url:"status,omitempty"`
	IdFrom        string `url:"id_from,omitempty"`
	IdTo          string `url:"id_to,omitempty"`
	CreatedFrom   string `url:"created_at_from,omitempty"`
	createdTo     string `url:"created_at_to,omitempty"`
}

func OdrHasExec(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.HasExecutions = value
		return nil
	}
}

func OrdSide(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.Side = value
		return nil
	}
}

func OrdSatus(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.Status = value
		return nil
	}
}

func OrdIdFrom(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.IdFrom = value
		return nil
	}
}

func OrdIdTo(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.IdTo = value
		return nil
	}
}

func OrdCreatedFrom(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.CreatedFrom = value
		return nil
	}
}

func OrdCreatedTo(value string) OrdersParams {
	return func(a *OrdersPameters) error {
		a.createdTo = value
		return nil
	}
}

func (a *Api) PlaceOrder(opts ...PlaceOrdersParams) models.CustomPlaceOrderInfo {
	orderInfo := models.CustomPlaceOrderInfo{}
	params := &PlaceOrdersPameters{}

	for _, op := range opts {
		err := op(params)
		if err != nil {
			return orderInfo
		}
	}

	order := models.PlaceOrderPayload{Async: true, Type: params.Type}

	price, _ := decimal.NewFromString(params.Price)
	pricestop, _ := decimal.NewFromString(params.PriceStop)

	cutPrice := strings.Split(price.String(), ".")
	limitPrice, _ := strconv.ParseInt(cutPrice[0], 10, 64)

	cutPriceStop, _ := strconv.ParseInt(pricestop.String(), 10, 64)

	switch params.Kind {
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

	order.Qty = params.Quantity

	c, err := caller.ClientWithToken(http.MethodPost, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientWithToken")
		}
		orderInfo.Error = err
		return orderInfo
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("ORDER_PLACE"),
		replacer.OptSymbol(params.Symbol),
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
		if config.Config.Debug {
			a.log.Error().Stack().Err(orderInfo.Error).Msg("Return orderInfo Error")
		}
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

	if err := a.CacheSetOrder([]models.OrdersIndex{
		{
			Symbol: params.Symbol,
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

func (a *Api) CancelAllCachedOrders(symbol string) error {

	if a.cache.GetDriver() == kind.GOCACHE {
		return fmt.Errorf("sorry, this method is not supported for GOCACHE driver")
	}

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

	if err := a.CacheSetOrder(ordersIndex); err != nil {
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
