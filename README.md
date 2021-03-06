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

### REST API

```go
package main

import (
	"fmt"
	"log"

	"github.com/cloudingcity/go-ftx/ftx"
)

func main() {
	client := ftx.New(
		ftx.WithAuth("your-api-key", "your-api-secret"),
		ftx.WithSubaccount("your-subaccount"), // Omit if not using subaccounts
	)
	account, err := client.Accounts.GetInformation()
	if err != nil {
		log.Fatal()
	}
	fmt.Printf("%+v", account)
}
```

### Websocket

```go
package main

import (
	"fmt"
	"log"

	"github.com/cloudingcity/go-ftx/ftx"
	"github.com/cloudingcity/go-ftx/ftx/stream"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := ftx.New(
		ftx.WithAuth("your-api-key", "your-api-secret"),
		ftx.WithSubaccount("your-subaccount"), // Omit if not using subaccounts
	)
	conn, err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	
	conn.PingRegular(ctx, 50*time.Second) // Keep connection prevent read timeout

	// Ping
	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}

	// Public Channels
	if err := conn.Subscribe(stream.ChannelTicker, "BTC/USD"); err != nil {
		log.Fatal(err)
	}
	if err := conn.Subscribe(stream.ChannelTrades, "BTC/USD"); err != nil {
		log.Fatal(err)
	}

	// Private Channels
	if err := conn.Login(); err != nil {
		log.Fatal(err)
	}
	if err := conn.Subscribe(stream.ChannelFills); err != nil {
		log.Fatal(err)
	}
	if err := conn.Subscribe(stream.ChannelOrders); err != nil {
		log.Fatal(err)
	}

	for {
		resp, err := conn.Recv()
		if err != nil {
			log.Fatal(err)
			return
		}

		switch v := resp.(type) {
		case stream.General:
			fmt.Println("general:", v)
		case stream.Pong:
			fmt.Println("pong:", v)
		case stream.OrderBook:
			fmt.Println("orderbook:", v)
		case stream.Trade:
			fmt.Println("trade:", v)
		case stream.Ticker:
			fmt.Println("ticker:", v)
		case stream.Fills:
			fmt.Println("fills:", v)
		case stream.Orders:
			fmt.Println("orders:", v)
		case stream.Error:
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
    - [x] Fills
    - [x] Orders
