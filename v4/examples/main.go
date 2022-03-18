package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thiagozs/go-cache/v1/cache/drivers/kind"
	"github.com/thiagozs/go-cache/v1/cache/options"
	"github.com/thiagozs/go-mbsdk/v4/api"
	"github.com/thiagozs/go-mbsdk/v4/pkg/cache"
)

func main() {
	key := os.Getenv("MB_KEY")
	secret := os.Getenv("MB_SECRET")

	optsc := []options.Options{
		options.OptFolder("./settings"),
		options.OptFileName("cache.db"),
		options.OptTTL(3000),
		options.OptLogDebug(true),
		options.OptLogDisable(false),
	}

	c, err := cache.NewCache(kind.BUNTDB, optsc...)
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

	auth, acc, err := a.Login()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", auth)
	fmt.Printf("%+v\n", acc)

	if balances, err := a.GetBalances(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(balances)
	}

	if ticker, err := a.Tickers("BTC-BRL"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ticker)
	}

}
