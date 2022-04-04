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

type WalletDepOptions func(c *WalletDepParameters) error
type WalletCoinOptions func(c *WalletCoinParameters) error

type WalletDepParameters struct {
	Limit  string `url:"limit,omitempty"`
	Page   string `url:"page,omitempty"`
	From   string `url:"from,omitempty"`
	To     string `url:"to,omitempty"`
	Symbol string `url:"symbol,omitempty"`
}

func WalletDepPage(page string) WalletDepOptions {
	return func(c *WalletDepParameters) error {
		c.Page = page
		return nil
	}
}

func WalletDepLimit(limit string) WalletDepOptions {
	return func(c *WalletDepParameters) error {
		c.Limit = limit
		return nil
	}
}

func WalletDepTo(to string) WalletDepOptions {
	return func(c *WalletDepParameters) error {
		c.To = to
		return nil
	}
}

func WalletDepFrom(from string) WalletDepOptions {
	return func(c *WalletDepParameters) error {
		c.From = from
		return nil
	}
}

func WalletDepSymbol(symbol string) WalletDepOptions {
	return func(c *WalletDepParameters) error {
		c.Symbol = symbol
		return nil
	}
}

type WalletCoinParameters struct {
	AccountRef  int    `url:"account_ref,omitempty"`
	Address     string `url:"address,omitempty"`
	Description string `url:"description,omitempty"`
	Quantity    string `url:"quantity,omitempty"`
	Symbol      string `url:"symbol,omitempty"`
	TxFee       string `url:"tx_fee,omitempty"`
}

func WalletCoinAccRef(accRef int) WalletCoinOptions {
	return func(c *WalletCoinParameters) error {
		c.AccountRef = accRef
		return nil
	}
}

func WalletCoinAddr(addr string) WalletCoinOptions {
	return func(c *WalletCoinParameters) error {
		c.Address = addr
		return nil
	}
}

func WalletCoinDesc(desc string) WalletCoinOptions {
	return func(c *WalletCoinParameters) error {
		c.Description = desc
		return nil
	}
}

func WalletCoinQty(quantity string) WalletCoinOptions {
	return func(c *WalletCoinParameters) error {
		c.Quantity = quantity
		return nil
	}
}

func WalletCoinSymbol(symbol string) WalletCoinOptions {
	return func(c *WalletCoinParameters) error {
		c.Symbol = symbol
		return nil
	}
}

func WalletCoinTxFee(tx_fee string) WalletCoinOptions {
	return func(c *WalletCoinParameters) error {
		c.TxFee = tx_fee
		return nil
	}
}

func (a *Api) WalletGetDeposit(opts ...WalletDepOptions) (models.WalletGetDepositsResponse, error) {
	deposits := models.WalletGetDepositsResponse{}
	params := &WalletDepParameters{}
	errApi := models.ErrorApiResponse{}

	for _, op := range opts {
		err := op(params)
		if err != nil {
			return deposits, err
		}
	}

	if params.Symbol == "" {
		return deposits, fmt.Errorf("symbol is required")
	}

	c, err := caller.ClientWithToken(http.MethodGet, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientWithToken")
		}
		return deposits, err
	}

	v, _ := query.Values(params)
	endpoint, err := replacer.Endpoint(
		replacer.OptKey("WALLET_DEPOSIT"),
		replacer.OptSymbol(params.Symbol),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return deposits, err
	}

	if (&WalletDepParameters{}) != params {
		endpoint = fmt.Sprintf("%s?%s", endpoint, v.Encode())
	}

	res, err := c.GetWithResponse(endpoint)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Get")
		}
		return deposits, err
	}
	defer res.Body.Close()

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ReadAll")
		}
		return deposits, err
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
			return deposits, err
		}
		return deposits, fmt.Errorf("%s - %s", errApi.Code, errApi.Message)
	}

	if err := json.Unmarshal(bts, &deposits); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal deposits")
		}
		return deposits, err
	}

	return deposits, nil
}

func (a *Api) WalletGetWithdrawCoin(symbol, withdrawId string) (models.WalletGetDepositsResponse, error) {
	withdrawcoin := models.WalletGetDepositsResponse{}
	errApi := models.ErrorApiResponse{}

	c, err := caller.ClientWithToken(http.MethodGet, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientWithToken")
		}
		return withdrawcoin, err
	}

	endpoint, err := replacer.Endpoint(
		replacer.OptKey("WALLET_GETWITHDRAW"),
		replacer.OptSymbol(symbol),
		replacer.OptWithDrawId(withdrawId),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return withdrawcoin, err
	}

	res, err := c.GetWithResponse(endpoint)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Get")
		}
		return withdrawcoin, err
	}
	defer res.Body.Close()

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ReadAll")
		}
		return withdrawcoin, err
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
			return withdrawcoin, err
		}
		return withdrawcoin, fmt.Errorf("%s - %s", errApi.Code, errApi.Message)
	}

	if err := json.Unmarshal(bts, &withdrawcoin); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal withdrawcoin")
		}
		return withdrawcoin, err
	}

	return withdrawcoin, nil
}

func (a *Api) WalletWithdrawCoin(opts ...WalletCoinOptions) (models.WalletWithdrawCoinResponse, error) {
	withdrawcoin := models.WalletWithdrawCoinResponse{}
	params := &WalletCoinParameters{}
	errApi := models.ErrorApiResponse{}

	for _, op := range opts {
		err := op(params)
		if err != nil {
			return withdrawcoin, err
		}
	}

	if params.Symbol == "" {
		return withdrawcoin, fmt.Errorf("symbol is required")
	}

	c, err := caller.ClientWithToken(http.MethodPost, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientWithToken")
		}
		return withdrawcoin, err
	}

	endpoint, err := replacer.Endpoint(
		replacer.OptKey("WALLET_GETDRAW"),
		replacer.OptSymbol(params.Symbol),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Replacer")
		}
		return withdrawcoin, err
	}

	wcp := models.WalletWithdrawCoinPayload{
		AccountRef:  params.AccountRef,
		Address:     params.Address,
		Description: params.Description,
		Quantity:    params.Quantity,
		Symbol:      params.Symbol,
		TxFee:       params.TxFee,
	}

	res, err := c.PostWithResponse(endpoint, wcp.ToBytes())
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("PostWithResponse")
		}
		return withdrawcoin, err
	}
	defer res.Body.Close()

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ReadAll")
		}
		return withdrawcoin, err
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
			return withdrawcoin, err
		}
		return withdrawcoin, fmt.Errorf("%s - %s", errApi.Code, errApi.Message)
	}

	if err := json.Unmarshal(bts, &withdrawcoin); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal withdrawcoin")
		}
		return withdrawcoin, err
	}

	return withdrawcoin, nil
}
