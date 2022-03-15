package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
	config.Config.Login = mts.key
	config.Config.Password = mts.secret
	config.Config.Cache = mts.cache
	config.Config.Debug = mts.debug
	config.Config.Endpoint = mts.endpoint

	return &Api{mts.cache}, nil
}

func (a *Api) AuthorizationToken() (models.AuthoritionToken, error) {
	auth := models.AuthoritionToken{}
	c, err := caller.ClientWithForm(http.MethodPost, a.cache)
	if err != nil {
		return auth, err
	}

	endpoint, err := replacer.Endpoint(replacer.OptKey("AUTHORIZE"),
		replacer.OptCache(a.cache),
	)
	if err != nil {
		return auth, err
	}

	res, err := c.PostFormWithResponse(endpoint)
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()
	if err != nil {
		return auth, err
	}

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return auth, err
	}

	if err := json.Unmarshal(bts, &auth); err != nil {
		return auth, err
	}

	if err := a.SetAuthorize(auth); err != nil {
		return auth, err
	}

	return auth, nil
}
