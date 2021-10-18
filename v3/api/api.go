package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/thiagozs/go-mbsdk/v3/models"
)

type Typo int

const (
	Ticker Typo = iota
	OrderBook
	Trades
	DaySummary
)

func (m Typo) String() string {
	return [...]string{
		"ticker",
		"orderbook",
		"trades",
		"day-summary",
	}[m]
}

var (
	api_base string = "https://www.mercadobitcoin.net/api"
	path     string = "{url}/{coin}/{method}/"
)

type Params struct {
	ExternalUrl  string   `json:"external"`
	Coin         string   `json:"coin"`
	Path         string   `json:"path"`
	Methods      *Methods `json:"methods"`
	HttpRetryMax int      `json:"http_maxretry"`
}

type Api struct {
	ExternalUrl  string
	Coin         string
	Url          string
	Methods      *Methods
	HttpRetryMax int
}

func New(params Params) *Api {
	var urlBuilder string

	res0 := strings.Replace(path, "{url}", api_base, -1)
	if len(params.ExternalUrl) > 0 {
		res0 = strings.Replace(api_base, "{url}", params.ExternalUrl, -1)
	}
	res1 := strings.Replace(res0, "{coin}", params.Coin, -1)
	urlBuilder = strings.Replace(res1, "{method}", params.Methods.Typo().String(), -1)

	switch params.Methods.Typo() {
	case Trades:
		if params.Methods.From() > 0 && params.Methods.To() == 0 {
			urlBuilder = fmt.Sprintf("%s%d/", urlBuilder, params.Methods.From())
		} else if params.Methods.From() > 0 && params.Methods.To() > 0 {
			urlBuilder = fmt.Sprintf("%s%d/%d/", urlBuilder, params.Methods.From(), params.Methods.To())
		}
	case DaySummary:
		urlBuilder = fmt.Sprintf("%s%d/%d/%d/", urlBuilder, params.Methods.Year(), params.Methods.Month(), params.Methods.Day())
	}

	return &Api{
		ExternalUrl:  params.ExternalUrl,
		Coin:         params.Coin,
		Url:          urlBuilder,
		Methods:      params.Methods,
		HttpRetryMax: params.HttpRetryMax,
	}
}

func (u *Api) GetURL() string {
	return u.Url
}

func (u *Api) GetCoin() string {
	return u.Coin
}

func (u *Api) FetchTiker() (models.TickerResponse, error) {
	response := models.TickerResponse{}
	buf, err := u.CallAPI()
	if err != nil {
		return response, err
	}

	if strings.Contains(string(buf[:]), "error") ||
		strings.Contains(string(buf[:]), "error_message") {
		return response, fmt.Errorf(string(buf[:]))
	}

	if err := json.Unmarshal(buf, &response); err != nil {
		return response, err
	}

	return response, nil
}

func (u *Api) FetchOrderBook() (models.OrderBookResponse, error) {
	response := models.OrderBookResponse{}
	buf, err := u.CallAPI()
	if err != nil {
		return response, err
	}

	if strings.Contains(string(buf[:]), "error") ||
		strings.Contains(string(buf[:]), "error_message") {
		return response, fmt.Errorf(string(buf[:]))
	}

	if err := json.Unmarshal(buf, &response); err != nil {
		return response, err
	}

	return response, nil
}

func (u *Api) FetchTrades() (models.TradesResponse, error) {
	response := models.TradesResponse{}
	buf, err := u.CallAPI()
	if err != nil {
		return response, err
	}

	if strings.Contains(string(buf[:]), "error") ||
		strings.Contains(string(buf[:]), "error_message") {
		return response, fmt.Errorf(string(buf[:]))
	}

	if err := json.Unmarshal(buf, &response); err != nil {
		return response, err
	}

	return response, nil
}

func (u *Api) DaySummary() (models.DaySummaryResponse, error) {
	response := models.DaySummaryResponse{}
	buf, err := u.CallAPI()
	if err != nil {
		return response, err
	}

	if strings.Contains(string(buf[:]), "error") ||
		strings.Contains(string(buf[:]), "error_message") {
		return response, fmt.Errorf(string(buf[:]))
	}

	if err := json.Unmarshal(buf, &response); err != nil {
		return response, err
	}

	return response, nil
}

func (u *Api) CallAPI() ([]byte, error) {
	client := retryablehttp.NewClient()
	client.RetryWaitMin = 10 * time.Millisecond
	client.RetryWaitMax = 10 * time.Millisecond
	client.RetryMax = u.HttpRetryMax

	// Create the request
	req, err := retryablehttp.NewRequest(http.MethodGet, u.GetURL(), nil)
	if err != nil {
		return []byte{}, err
	}

	// Send the request.
	rr, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	if err != nil {
		return []byte(""), err
	}
	defer rr.Body.Close()

	return ioutil.ReadAll(rr.Body)
}
