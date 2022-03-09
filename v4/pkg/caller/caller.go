package caller

import (
	"encoding/json"

	"github.com/thiagozs/go-mbsdk/v4/config"
	"github.com/thiagozs/go-mbsdk/v4/models"
	"github.com/thiagozs/go-mbsdk/v4/pkg/cache"
	"github.com/thiagozs/go-mbsdk/v4/pkg/client"
)

func ClientWithToken(method string, g *cache.Cache) (client.HttpClientPort, error) {
	c := client.NewHttpClient(3, 3, 3)
	c.DisableLogLevel()
	c.SetHeader(method, "Content-Type", "application/json")
	c.SetHeader(method, "Accept", "*/*")
	c.SetHeader(
		method,
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36",
	)
	jraw, err := g.GetKeyVal(config.AUTHORIZE.String())
	if err != nil {
		return c, err
	}

	auth := models.AuthoritionToken{}
	if err := json.Unmarshal([]byte(jraw), &auth); err != nil {
		return c, err
	}

	c.SetHeader(method, "Authorization", "Bearer "+auth.AccessToken)
	return c, nil
}

func ClientPublic(method string, g *cache.Cache) (client.HttpClientPort, error) {
	c := client.NewHttpClient(3, 3, 3)
	c.DisableLogLevel()
	c.SetHeader(method, "Content-Type", "application/json")
	c.SetHeader(method, "Accept", "*/*")
	c.SetHeader(
		method,
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36",
	)
	return c, nil
}

func ClientWithForm(method string, g *cache.Cache) (client.HttpClientPort, error) {
	c := client.NewHttpClient(3, 3, 3)
	c.DisableLogLevel()
	c.SetHeader(method, "Content-Type", "application/x-www-form-urlencoded")
	c.SetHeader(method, "Accept", "*/*")
	c.SetHeader(
		method,
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36",
	)

	c.SetFormValue(method, "login", config.Config.Login)
	c.SetFormValue(method, "password", config.Config.Password)
	return c, nil
}
