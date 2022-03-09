package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thiagozs/go-mbsdk/v4/api"
	"github.com/thiagozs/go-mbsdk/v4/pkg/cache"
)

func main() {
	key := os.Getenv("MB_KEY")
	secret := os.Getenv("MB_SECRET")

	c, err := cache.NewCache()
	if err != nil {
		log.Fatal(err)
	}

	opts := []api.Options{
		api.OptKey(key),
		api.OptSecret(secret),
		api.OptDebug(true),
		api.OptCache(c),
	}
	a, err := api.New(opts...)
	if err != nil {
		fmt.Println(err)
	}

	if auth, err := a.AuthorizationToken(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(auth)
	}

	if ticker, err := a.Tickers("BTC-BRL"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ticker)
	}

}
