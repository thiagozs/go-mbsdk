# MercadoBitcoin SDK

Easy way to consume the public api informations from MercadoBitcoin

## API v4 (new - working in progress)

### Endpoints available

- [x] Authorization
- [ ] Accounts
	- [x] - Get Accounts
	- [x] - Balance List
	- [ ] - Position List
- [x] Trading
	- [x] - Get Order
	- [x] - Order Place
	- [x] - Order Cancel
	- [x] - Order List
	- [x] - Order Cancel All
- [x] Wallet
	- [x] Wallet Deposit
	- [x] Wallet Withdraw
	- [x] Wallet GetDraw
- [x] Public data
	- [x] - Get Ticker
	- [x] - Get Orderbook
	- [x] - Get Trades
	- [x] - Get Candles
	- [x] - Get Symbol

### Cache
The external cache system is not mandatory, but if you want to use a functions worked with cache for a delayed cli command, you needed use the cache system.

```golang
key := os.Getenv("MB_KEY")
secret := os.Getenv("MB_SECRET")

// not mandatory
optsc := []options.Options{
	options.OptFolder("./settings"),
	options.OptFileName("cache.db"),
	options.OptTTL(3000),
	options.OptLogDebug(true),
	options.OptLogDisable(false),
}

// not mandatory
c, err := cache.NewCache(kind.BUNTDB, optsc...)
if err != nil {
	log.Fatal(err)
}

opts := []api.Options{
	api.OptKey(key),
	api.OptSecret(secret),
	api.OptDebug(true),
	api.OptCache(c), // not mandatory
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
```

## ~~Example of API consume of Version 3~~ (deprecated)

Simple code writed on `main.go`. Just fill up the model **params** and **methods** and make a request.

## Versioning and license

Our version numbers follow the [semantic versioning specification](http://semver.org/). You can see the available versions by checking the [tags on this repository](https://github.com/thiagozs/go-mbsdk/tags). For more details about our license model, please take a look at the [LICENSE](LICENSE) file.

2022, thiagozs
