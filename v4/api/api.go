package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/thiagozs/go-mbsdk/v4/config"
	"github.com/thiagozs/go-mbsdk/v4/models"
	"github.com/thiagozs/go-mbsdk/v4/pkg/caller"
	"github.com/thiagozs/go-mbsdk/v4/pkg/replacer"
)

func New(opts ...Options) (*Api, error) {

	mts := &ApiCfg{}

	for _, op := range opts {
		err := op(mts)
		if err != nil {
			return &Api{}, err
		}
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log := zerolog.New(os.Stderr).With().
		Caller().
		Timestamp().Logger()

	config.Config.Login = mts.key
	config.Config.Password = mts.secret
	config.Config.Cache = mts.cache
	config.Config.Debug = mts.debug
	config.Config.Endpoint = mts.endpoint

	return &Api{mts.cache, log}, nil
}

func (a *Api) AuthorizationToken() (models.AuthoritionToken, error) {
	auth := models.AuthoritionToken{}
	c, err := caller.ClientWithForm(http.MethodPost, a.cache)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ClientWithForm")
		}
		return auth, err
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("AUTHORIZE"),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("EndPoint")
		}
		return auth, err
	}

	res, err := c.PostFormWithResponse(endpoint)
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("PostFormWithResponse")
		}
		return auth, err
	}

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("ReadAll")
		}
		return auth, err
	}

	if err := json.Unmarshal(bts, &auth); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("Json Unmarshal AuthoritionToken")
		}
		return auth, err
	}

	if err := a.SetAuthorize(auth); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("SetAuthorize")
		}
		return auth, err
	}

	return auth, nil
}
