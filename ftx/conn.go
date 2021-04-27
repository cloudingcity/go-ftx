package ftx

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

type WSCommon struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Market  string `json:"market"`
}

type WSError struct {
	Type string `json:"type"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type WSPong struct {
	Type string `json:"op"`
}

type WSOrderBook struct {
	WSCommon
	Data struct {
		Bids     [][]float64 `json:"bids"`
		Asks     [][]float64 `json:"asks"`
		Time     *Time       `json:"time"`
		Checksum int         `json:"checksum"`
		Action   string      `json:"action"`
	} `json:"data"`
}

type WSTrade struct {
	WSCommon
	Data []struct {
		Trade
	} `json:"data"`
}

type WSTicker struct {
	WSCommon
	Data struct {
		Bid     float64 `json:"bid"`
		Ask     float64 `json:"ask"`
		BidSize float64 `json:"bidSize"`
		AskSize float64 `json:"askSize"`
		Last    float64 `json:"last"`
		Time    *Time   `json:"time"`
	} `json:"data"`
}
type Conn struct {
	conn *websocket.Conn
}

func (c *Conn) Recv() (interface{}, error) {
	var resp ConnResponse
	if err := c.conn.ReadJSON(&resp); err != nil {
		return nil, err
	}

	if resp.Type == "error" {
		return WSError{Type: resp.Type, Code: resp.Code, Msg: resp.Msg}, nil
	}

	if resp.Type == "pong" {
		return WSPong{Type: resp.Type}, nil
	}

	if resp.Data == nil {
		return WSCommon{Type: resp.Type, Channel: resp.Channel, Market: resp.Market}, nil
	}

	switch resp.Channel {
	case "orderbook":
		v := WSOrderBook{WSCommon: WSCommon{Type: resp.Type, Channel: resp.Channel, Market: resp.Market}}
		err := json.Unmarshal(resp.Data, &v.Data)
		return v, err
	case "trades":
		v := WSTrade{WSCommon: WSCommon{Type: resp.Type, Channel: resp.Channel, Market: resp.Market}}
		err := json.Unmarshal(resp.Data, &v.Data)
		return v, err
	case "ticker":
		v := WSTicker{WSCommon: WSCommon{Type: resp.Type, Channel: resp.Channel, Market: resp.Market}}
		err := json.Unmarshal(resp.Data, &v.Data)
		return v, err
	default:
		return nil, fmt.Errorf("channel %q not support", resp.Channel)
	}
}

func (c *Conn) Ping() error {
	return c.conn.WriteJSON(&ConnRequest{OP: "ping"})
}

func (c *Conn) Subscribe(channel, market string) error {
	return c.conn.WriteJSON(
		&ConnRequest{
			OP:      "subscribe",
			Channel: channel,
			Market:  market,
		},
	)
}

func (c *Conn) Unsubscribe(channel, market string) error {
	return c.conn.WriteJSON(
		&ConnRequest{
			OP:      "unsubscribe",
			Channel: channel,
			Market:  market,
		},
	)
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
