package ftx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestMarketService_All(t *testing.T) {
	c, srv, teardown := setup()
	defer teardown()

	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.SetBodyString(`{"success":true,"result":[{"name":"BTC/USD"}]}`)
	}

	markets, err := c.Markets.All()

	assert.NoError(t, err)
	assert.Equal(t, "BTC/USD", markets[0].Name)
}

func TestMarketService_Get(t *testing.T) {
	c, srv, teardown := setup()
	defer teardown()

	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.SetBodyString(`{"success":true,"result":{"name":"BTC/USD"}}`)
	}

	market, err := c.Markets.Get("BTC/USD")

	assert.NoError(t, err)
	assert.Equal(t, "BTC/USD", market.Name)
}

func TestMarketService_GetOrderBook(t *testing.T) {
	c, srv, teardown := setup()
	defer teardown()

	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.SetBodyString(`{"success":true,"result":{"asks":[[111,222]],"bids":[[333,444]]}}`)
	}

	orderbook, err := c.Markets.GetOrderBook("BTC/USD", nil)

	assert.NoError(t, err)
	assert.Equal(t, float64(111), orderbook.Asks[0][0])
	assert.Equal(t, float64(333), orderbook.Bids[0][0])
}

func TestMarketService_GetTrades(t *testing.T) {
	c, srv, teardown := setup()
	defer teardown()

	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.SetBodyString(`{"success":true,"result":[{"id":123456}]}`)
	}

	trades, err := c.Markets.GetTrades("BTC/USD", nil)

	assert.NoError(t, err)
	assert.Equal(t, 123456, trades[0].Id)
}

func TestMarketService_GetHistoricalPrices(t *testing.T) {
	c, srv, teardown := setup()
	defer teardown()

	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.SetBodyString(`{"success":true,"result":[{"open":123456}]}`)
	}

	candles, err := c.Markets.GetHistoricalPrices("BTC/USD", nil)

	assert.NoError(t, err)
	assert.Equal(t, float64(123456), candles[0].Open)
}
