# MercadoBitcoin SDK

Easy way to consume the public api informations from MercadoBitcoin

## API v4 (new)

```golang
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
```

## ~~Example of API consume of Version 3~~ (deprecated)

Simple code writed on `main.go`. Just fill up the model **params** and **methods** and make a request.

```golang
package main

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
```
## Versioning and license

Our version numbers follow the [semantic versioning specification](http://semver.org/). You can see the available versions by checking the [tags on this repository](https://github.com/thiagozs/go-mbsdk/tags). For more details about our license model, please take a look at the [LICENSE](LICENSE) file.

2022, thiagozs