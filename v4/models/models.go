package models

import "encoding/json"

type AuthoritionToken struct {
	AccessToken string `json:"access_token"`
	Expiration  int    `json:"expiration"`
}

func (p *AuthoritionToken) ToBytes() []byte {
	bts, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return bts
}

type TickersQuery struct {
	Symbols string `url:"symbols"`
}

type TickersResponse []struct {
	Buy  string `json:"buy"`
	Date int    `json:"date"`
	High string `json:"high"`
	Last string `json:"last"`
	Low  string `json:"low"`
	Open string `json:"open"`
	Pair string `json:"pair"`
	Sell string `json:"sell"`
	Vol  string `json:"vol"`
}

type ListBalancesResponse []struct {
	Available string `json:"available"`
	OnHold    string `json:"on_hold"`
	Symbol    string `json:"symbol"`
	Total     string `json:"total"`
}

type PlaceOrderPayload struct {
	Async      bool   `json:"async,omitempty"`
	Cost       int    `json:"cost,omitempty"`
	LimitPrice int    `json:"limitPrice,omitempty"`
	Qty        string `json:"qty,omitempty"`
	Side       string `json:"side,omitempty"`
	StopPrice  int    `json:"stopPrice,omitempty"`
	Type       string `json:"type,omitempty"`
}

func (p *PlaceOrderPayload) ToBytes() []byte {
	bts, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return bts
}

type PlaceOrderResponse struct {
	OrderID string `json:"orderId"`
}

func (p *PlaceOrderResponse) ToBytes() []byte {
	bts, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return bts
}

type ListAccountsResponse []struct {
	Currency     string `json:"currency"`
	CurrencySign string `json:"currencySign"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}

type ErrorPlaceOrderResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (p *ErrorPlaceOrderResponse) ToBytes() []byte {
	bts, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return bts
}

type CustomPlaceOrderInfo struct {
	StatusCode int    `json:"status_code"`
	OrderID    string `json:"order_id"`
	Payload    string `json:"payload"`
	Response   string `json:"response"`
	EndPoint   string `json:"endpoint"`
	Error      error  `json:"error"`
}

type ListPositionResponse []struct {
	AvgPrice   int    `json:"avgPrice"`
	Category   string `json:"category"`
	ID         string `json:"id"`
	Instrument string `json:"instrument"`
	Qty        string `json:"qty"`
	Side       string `json:"side"`
}

type ListOrderResponse []GetOrderResponse

type GetOrderResponse struct {
	AvgPrice   int `json:"avgPrice"`
	CreatedAt  int `json:"created_at"`
	Executions []struct {
		ExecutedAt int    `json:"executed_at"`
		FeeRate    string `json:"fee_rate"`
		ID         string `json:"id"`
		Instrument string `json:"instrument"`
		Price      int    `json:"price"`
		Qty        string `json:"qty"`
		Side       string `json:"side"`
	} `json:"executions"`
	Fee            string `json:"fee"`
	FilledQty      string `json:"filledQty"`
	ID             string `json:"id"`
	Instrument     string `json:"instrument"`
	LimitPrice     int    `json:"limitPrice"`
	Qty            string `json:"qty"`
	Side           string `json:"side"`
	Status         string `json:"status"`
	StopPrice      int    `json:"stopPrice"`
	TriggerOrderID string `json:"triggerOrderId"`
	Type           string `json:"type"`
	UpdatedAt      int    `json:"updated_at"`
}

type OrdersIndex struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Side   string `json:"side"`
	Type   string `json:"type"`
	Price  string `json:"price"`
}

type OrdersIndexResponse []OrdersIndex
