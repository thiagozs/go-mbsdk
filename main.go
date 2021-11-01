package main

import (
	"fmt"
	"log"

	"github.com/thiagozs/go-mbsdk/v3/api"
)

func main() {

	req := api.New(api.Params{
		Coin: "BTC",
		Options: []api.MethodsOpts{
			api.OptsType(api.Ticker),
		},
		HttpRetryMax: 3,
	})

	res, err := req.FetchTiker()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", res)
}
