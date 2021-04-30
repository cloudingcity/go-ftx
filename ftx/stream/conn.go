package stream

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

const (
	ChannelOrderBook = "orderbook"
	ChannelTrades    = "trades"
	ChannelTicker    = "ticker"
)

type ConnRequest struct {
	OP      string `json:"op"`
	Channel string `json:"channel,omitempty"`
	Market  string `json:"market,omitempty"`
}

type ConnResponse struct {
	Type    string          `json:"type"`
	Channel string          `json:"channel,omitempty"`
	Market  string          `json:"market,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Code    int             `json:"code,omitempty"`
	Msg     string          `json:"msg,omitempty"`
}

type Conn struct {
	Conn *websocket.Conn
}

func (c *Conn) Recv() (interface{}, error) {
	var resp ConnResponse
	if err := c.Conn.ReadJSON(&resp); err != nil {
		return nil, err
	}

	if resp.Type == "error" {
		return Error{Type: resp.Type, Code: resp.Code, Msg: resp.Msg}, nil
	}

	if resp.Type == "pong" {
		return Pong{Type: resp.Type}, nil
	}

	if resp.Data == nil {
		return General{Type: resp.Type, Channel: resp.Channel, Market: resp.Market}, nil
	}

	switch resp.Channel {
	case "orderbook":
		v := OrderBook{General: General{Type: resp.Type, Channel: resp.Channel, Market: resp.Market}}
		err := json.Unmarshal(resp.Data, &v.Data)
		return v, err
	case "trades":
		v := Trade{General: General{Type: resp.Type, Channel: resp.Channel, Market: resp.Market}}
		err := json.Unmarshal(resp.Data, &v.Data)
		return v, err
	case "ticker":
		v := Ticker{General: General{Type: resp.Type, Channel: resp.Channel, Market: resp.Market}}
		err := json.Unmarshal(resp.Data, &v.Data)
		return v, err
	default:
		return nil, fmt.Errorf("channel %q not support", resp.Channel)
	}
}

func (c *Conn) Ping() error {
	return c.Conn.WriteJSON(&ConnRequest{OP: "ping"})
}

func (c *Conn) Subscribe(channel, market string) error {
	return c.Conn.WriteJSON(
		&ConnRequest{
			OP:      "subscribe",
			Channel: channel,
			Market:  market,
		},
	)
}

func (c *Conn) Unsubscribe(channel, market string) error {
	return c.Conn.WriteJSON(
		&ConnRequest{
			OP:      "unsubscribe",
			Channel: channel,
			Market:  market,
		},
	)
}

func (c *Conn) Close() error {
	return c.Conn.Close()
}
