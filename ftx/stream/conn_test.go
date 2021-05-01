package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConn_Ping(t *testing.T) {
	conn, _, teardown := setup()
	defer teardown()

	err := conn.Ping()
	assert.NoError(t, err)

	resp, err := conn.RecvRaw()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"op":"ping"}`, string(resp))
}

// examples from https://docs.ftx.com/#authentication-2
func TestConn_Login(t *testing.T) {
	conn, _, teardown := setup()
	defer teardown()

	const (
		key    = "api-key"
		secret = "Y2QTHI23f23f23jfjas23f23To0RfUwX3H42fvN-"
	)
	unixTime = func() int64 { return 1557246346499 }
	conn.key = key
	conn.secret = []byte(secret)

	err := conn.Login()
	assert.NoError(t, err)

	resp, err := conn.RecvRaw()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"op":"login", "args":{"key":"api-key","sign":"d10b5a67a1a941ae9463a60b285ae845cdeac1b11edc7da9977bef0228b96de9","time":1557246346499}}`, string(resp))
}

func TestConn_Subscribe(t *testing.T) {
	conn, _, teardown := setup()
	defer teardown()

	t.Run("public channel", func(t *testing.T) {
		err := conn.Subscribe(ChannelTrades, "BTC/USD")
		assert.NoError(t, err)

		resp, err := conn.RecvRaw()

		assert.NoError(t, err)
		assert.JSONEq(t, `{"op":"subscribe","channel":"trades","market":"BTC/USD"}`, string(resp))
	})

	t.Run("private channel", func(t *testing.T) {
		err := conn.Subscribe(ChannelFills)
		assert.NoError(t, err)

		resp, err := conn.RecvRaw()

		assert.NoError(t, err)
		assert.JSONEq(t, `{"op":"subscribe","channel":"fills"}`, string(resp))
	})
}

func TestConn_Unsubscribe(t *testing.T) {
	conn, _, teardown := setup()
	defer teardown()

	t.Run("public channel", func(t *testing.T) {
		err := conn.Unsubscribe(ChannelTrades, "BTC/USD")
		assert.NoError(t, err)

		resp, err := conn.RecvRaw()

		assert.NoError(t, err)
		assert.JSONEq(t, `{"op":"unsubscribe","channel":"trades","market":"BTC/USD"}`, string(resp))
	})

	t.Run("private channel", func(t *testing.T) {
		err := conn.Unsubscribe(ChannelFills)
		assert.NoError(t, err)

		resp, err := conn.RecvRaw()

		assert.NoError(t, err)
		assert.JSONEq(t, `{"op":"unsubscribe","channel":"fills"}`, string(resp))
	})
}

func TestConn_Recv(t *testing.T) {
	conn, ws, teardown := setup()
	defer teardown()

	t.Run("error", func(t *testing.T) {
		_ = ws.WriteJSON(&connResponse{Type: "error", Code: 500, Msg: "something wrong"})

		resp, err := conn.Recv()
		assert.NoError(t, err)
		assert.IsType(t, Error{}, resp)

		got := resp.(Error)
		assert.Equal(t, "error", got.Type)
		assert.Equal(t, 500, got.Code)
		assert.Equal(t, "something wrong", got.Msg)
	})

	t.Run("pong", func(t *testing.T) {
		_ = ws.WriteJSON(&connResponse{Type: "pong"})

		resp, err := conn.Recv()
		assert.NoError(t, err)
		assert.IsType(t, Pong{}, resp)

		got := resp.(Pong)
		assert.Equal(t, "pong", got.Type)
	})

	t.Run("general", func(t *testing.T) {
		_ = ws.WriteJSON(&connResponse{Type: "subscribed", Channel: ChannelTicker, Market: "BTC/USD"})

		resp, err := conn.Recv()
		assert.NoError(t, err)
		assert.IsType(t, General{}, resp)

		got := resp.(General)
		assert.Equal(t, "subscribed", got.Type)
		assert.Equal(t, ChannelTicker, got.Channel)
		assert.Equal(t, "BTC/USD", got.Market)
	})

	t.Run("orderbook", func(t *testing.T) {
		_ = ws.WriteJSON(&connResponse{
			Type:    "update",
			Channel: ChannelOrderBook,
			Market:  "BTC/USD",
			Data:    []byte(`{"action":"buy"}`),
		})

		resp, err := conn.Recv()
		assert.NoError(t, err)
		assert.IsType(t, OrderBook{}, resp)

		got := resp.(OrderBook)
		assert.Equal(t, "buy", got.Data.Action)
	})

	t.Run("trades", func(t *testing.T) {
		_ = ws.WriteJSON(&connResponse{
			Type:    "update",
			Channel: ChannelTrades,
			Market:  "BTC/USD",
			Data:    []byte(`[{"id":123}]`),
		})

		resp, err := conn.Recv()
		assert.NoError(t, err)
		assert.IsType(t, Trade{}, resp)

		got := resp.(Trade)
		assert.Equal(t, 123, got.Data[0].ID)
	})

	t.Run("ticker", func(t *testing.T) {
		_ = ws.WriteJSON(&connResponse{
			Type:    "update",
			Channel: ChannelTicker,
			Market:  "BTC/USD",
			Data:    []byte(`{"last":1234.56}`),
		})

		resp, err := conn.Recv()
		assert.NoError(t, err)
		assert.IsType(t, Ticker{}, resp)

		got := resp.(Ticker)
		assert.Equal(t, 1234.56, got.Data.Last)
	})

	t.Run("fills", func(t *testing.T) {
		_ = ws.WriteJSON(&connResponse{
			Type:    "update",
			Channel: ChannelFills,
			Market:  "BTC/USD",
			Data:    []byte(`{"id":123}`),
		})

		resp, err := conn.Recv()
		assert.NoError(t, err)
		assert.IsType(t, Fills{}, resp)

		got := resp.(Fills)
		assert.Equal(t, 123, got.Data.ID)
	})

	t.Run("orders", func(t *testing.T) {
		_ = ws.WriteJSON(&connResponse{
			Type:    "update",
			Channel: ChannelOrders,
			Market:  "BTC/USD",
			Data:    []byte(`{"id":123}`),
		})

		resp, err := conn.Recv()
		assert.NoError(t, err)
		assert.IsType(t, Orders{}, resp)

		got := resp.(Orders)
		assert.Equal(t, 123, got.Data.ID)
	})
}
