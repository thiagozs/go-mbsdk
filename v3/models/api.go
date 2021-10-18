package models

import "encoding/json"

type DaySummaryResponse struct {
	Date     string      `json:"date,omitempty"`
	Opening  json.Number `json:"opening,omitempty"`
	Closing  json.Number `json:"closing,omitempty"`
	Lowest   json.Number `json:"lowest,omitempty"`
	Highest  json.Number `json:"highest,omitempty"`
	Volume   json.Number `json:"volume,omitempty"`
	Quantity json.Number `json:"quantity,omitempty"`
	Amount   json.Number `json:"amount,omitempty"`
	AvgPrice json.Number `json:"avg_price,omitempty"`
}

type OrderBookResponse struct {
	Asks [][]json.Number `json:"asks,omitempty"`
	Bids [][]json.Number `json:"bids,omitempty"`
}

type TickerResponse struct {
	Ticker struct {
		High json.Number `json:"high,omitempty"`
		Low  json.Number `json:"low,omitempty"`
		Vol  json.Number `json:"vol,omitempty"`
		Last json.Number `json:"last,omitempty"`
		Buy  json.Number `json:"buy,omitempty"`
		Sell json.Number `json:"sell,omitempty"`
		Date json.Number `json:"date,omitempty"`
	} `json:"ticker"`
}

type TradesResponse []struct {
	Tid    json.Number `json:"tid,omitempty"`
	Date   json.Number `json:"date,omitempty"`
	Type   string      `json:"type,omitempty"`
	Price  json.Number `json:"price,omitempty"`
	Amount json.Number `json:"amount,omitempty"`
}
