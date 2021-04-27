# go-ftx

[![Test](https://github.com/cloudingcity/go-ftx/workflows/Test/badge.svg)](https://github.com/cloudingcity/go-ftx/actions?query=workflow%3ATest)
[![Lint](https://github.com/cloudingcity/go-ftx/workflows/Lint/badge.svg)](https://github.com/cloudingcity/go-ftx/actions?query=workflow%3ALint)
[![codecov](https://codecov.io/gh/cloudingcity/go-ftx/branch/main/graph/badge.svg)](https://codecov.io/gh/cloudingcity/go-ftx)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudingcity/go-ftx)](https://goreportcard.com/report/github.com/cloudingcity/go-ftx)

go-ftx is a Go client library for accessing the [FTX API](https://docs.ftx.com/).

# Install

```console
go get github.com/cloudingcity/go-ftx
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/cloudingcity/go-ftx/ftx"
)

func main() {
	client := ftx.New()
	market, err := client.Markets.Get("ETH/USD")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v", market)
}
```

## Examples

### Get Private Resource

```go
client := ftx.New(
    ftx.WithAuth("your-api-key", "your-api-secret"),
    ftx.WithSubAccount("your-subaccount"), // Omit if not using subaccounts
)
account, err := client.Accounts.GetInformation()
```

### Websocket

```go
package main

import (
	"fmt"
	"log"

	"github.com/cloudingcity/go-ftx/ftx"
)

func main() {
	c := ftx.New()
	conn, err := c.Connect()
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}
	if err := conn.Subscribe(ftx.ChannelTicker, "BTC/USD"); err != nil {
		log.Fatal(err)
	}

	for {
		resp, err := conn.Recv()
		if err != nil {
			log.Fatal(err)
			return
		}

		switch v := resp.(type) {
		case ftx.WSCommon:
			fmt.Println("common:", v)
		case ftx.WSPong:
			fmt.Println("pong:", v)
		case ftx.WSOrderBook:
			fmt.Println("orderbook:", v)
		case ftx.WSTrade:
			fmt.Println("trade:", v)
		case ftx.WSTicker:
			fmt.Println("ticker:", v)
		case ftx.WSError:
			fmt.Println("error:", v)
		}
	}
}
```

## Todos

- [ ] REST API
    - [x] Marktes
    - [x] Accounts
    - [ ] Subaccounts
    - [ ] Futures
    - [ ] Wallet
    - [ ] Orders
    - [ ] Convert
    - [ ] Spot Margin
    - [ ] Fills
    - [ ] Funding Payments
    - [ ] Leveraged Tokens
    - [ ] Options
    - [ ] Staking
- [ ] Websocket API
    - [x] Ping
    - [x] OrderBooks
    - [x] Trade
    - [x] Ticker
    - [ ] Markets
    - [ ] Grouped Orderbooks
    - [ ] Fills
    - [ ] Orders
