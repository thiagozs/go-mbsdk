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

	c, err := cache.NewCache(false, true)
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

	// step 1 - get all authorization (mandatory)
	if auth, err := a.AuthorizationToken(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(auth)
	}

	// step 2 - run the account function to get all info (mandatory)
	if acc, err := a.GetAccounts(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(acc)
	}

	// step 3 - run others methods, before that, you need
	// run the function account to get all information
	if balances, err := a.GetBalances(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(balances)
	}

	// step 4 - publics endpoints
	if ticker, err := a.Tickers("BTC-BRL"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ticker)
	}

}
