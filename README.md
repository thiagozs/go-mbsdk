# MercadoBitcoin SDK

Easy way to consume the public api informations from MercadoBitcoin

## Example of API consume of Version 3

Simple code writed on `main.go`. Just fill up the model **params** and **methods** and make a request.

```golang
package main

func main() {
	req := api.New(api.Params{
		Coin:         "BTC",
		Methods:      api.NewMethods(0, 0, 0, 0, 0, api.Ticker),
		HttpRetryMax: 3,
	})

	res, err := req.FetchTiker()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", res)
}
```