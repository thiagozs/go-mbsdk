package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/thiagozs/go-cache/v1/cache/drivers/kind"
	"github.com/thiagozs/go-cache/v1/cache/options"
	"github.com/thiagozs/go-mbsdk/v4/config"
	"github.com/thiagozs/go-mbsdk/v4/models"
	"github.com/thiagozs/go-mbsdk/v4/pkg/cache"
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

	if mts.cache == nil {
		cache, err := cache.NewCache(kind.GOCACHE,
			options.OptTimeCleanUpInt(time.Duration(60)*time.Second),
			options.OptTimeExpiration(time.Duration(300)*time.Second))
		if err != nil {
			return &Api{}, err
		}
		mts.cache = cache
	}

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

	if err := a.CacheSetAuthorize(auth); err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("SetAuthorize")
		}
		return auth, err
	}

	return auth, nil
}

func (a *Api) Login() (models.AuthoritionToken, models.ListAccountsResponse, error) {
	auth, err := a.AuthorizationToken()
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("AuthorizationToken")
		}
		return models.AuthoritionToken{}, models.ListAccountsResponse{}, err
	}

	acc, err := a.GetAccounts()
	if err != nil {
		if config.Config.Debug {
			a.log.Error().Stack().Err(err).Msg("GetAccounts")
		}
		return models.AuthoritionToken{}, models.ListAccountsResponse{}, err
	}

	return auth, acc, nil
}
