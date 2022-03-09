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

type GetBalancesResponse []struct {
	Available json.Number `json:"available"`
	Symbol    string      `json:"symbol"`
	Total     json.Number `json:"total"`
}

type PlaceOrderPayload struct {
	Async      bool    `json:"async,omitempty"`
	Cost       int     `json:"cost,omitempty"`
	LimitPrice int     `json:"limitPrice,omitempty"`
	Qty        float64 `json:"qty,omitempty"`
	Side       string  `json:"side,omitempty"`
	Type       string  `json:"type,omitempty"`
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

type GetOrderResponse struct {
	AvgPrice   int `json:"avgPrice"`
	CreatedAt  int `json:"created_at"`
	Executions []struct {
		ExecutedAt int         `json:"executed_at"`
		ID         string      `json:"id"`
		Instrument string      `json:"instrument"`
		Price      int         `json:"price"`
		Qty        json.Number `json:"qty"`
		Side       string      `json:"side"`
	} `json:"executions"`
	FilledQty  json.Number `json:"filledQty"`
	ID         string      `json:"id"`
	Instrument string      `json:"instrument"`
	LimitPrice int         `json:"limitPrice"`
	Qty        json.Number `json:"qty"`
	Side       string      `json:"side"`
	Status     string      `json:"status"`
	Type       string      `json:"type"`
	UpdatedAt  int         `json:"updated_at"`
}

type GetAccountsResponse []struct {
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
