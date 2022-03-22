package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/thiagozs/go-cache/v1/cache/drivers/kind"
	"github.com/thiagozs/go-cache/v1/cache/options"
	"github.com/thiagozs/go-mbsdk/v4/api"
	"github.com/thiagozs/go-mbsdk/v4/pkg/cache"
)

func main() {
	key := os.Getenv("MB_KEY")
	secret := os.Getenv("MB_SECRET")
	endpoint := os.Getenv("MB_ENDPOINT")

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
		api.OptEndpoint(endpoint),
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

	if ticker, err := a.OrderBook("BTC-BRL", "1"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ticker)
	}

	if trades, err := a.Trades("BTC-BRL"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(trades)
	}

	if trades, err := a.Symbols([]string{"BTC-BRL"}); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(trades)
	}

	now := int(time.Now().Unix())

	if candles, err := a.Candles(api.CandSymbols("BTC-BRL"),
		api.CandResolution("15m"), api.CandFrom(now),
		api.CandTo(now+3600)); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(candles)
	}
}
