package stream

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	ChannelOrderBook = "orderbook"
	ChannelTrades    = "trades"
	ChannelTicker    = "ticker"
	ChannelFills     = "fills"
	ChannelOrders    = "orders"
)

type connRequest struct {
	OP      string `json:"op"`
	Args    *args  `json:"args,omitempty"`
	Channel string `json:"channel,omitempty"`
	Market  string `json:"market,omitempty"`
}

type args struct {
	Key        string `json:"key"`
	Sign       string `json:"sign"`
	Time       int64  `json:"time"`
	SubAccount string `json:"subaccount,omitempty"`
}

type connResponse struct {
	Type    string          `json:"type"`
	Channel string          `json:"channel,omitempty"`
	Market  string          `json:"market,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Code    int             `json:"code,omitempty"`
	Msg     string          `json:"msg,omitempty"`
}

type Conn struct {
	conn       *websocket.Conn
	key        string
	secret     []byte
	subaccount string
}

func New(conn *websocket.Conn, key string, secret []byte, subaccount string) *Conn {
	return &Conn{conn: conn, key: key, secret: secret, subaccount: subaccount}
}

func (c *Conn) Recv() (interface{}, error) {
	var resp connResponse
	if err := c.conn.ReadJSON(&resp); err != nil {
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
	case ChannelOrderBook:
		v := OrderBook{General: General{Type: resp.Type, Channel: resp.Channel, Market: resp.Market}}
		err := json.Unmarshal(resp.Data, &v.Data)
		return v, err
	case ChannelTrades:
		v := Trade{General: General{Type: resp.Type, Channel: resp.Channel, Market: resp.Market}}
		err := json.Unmarshal(resp.Data, &v.Data)
		return v, err
	case ChannelTicker:
		v := Ticker{General: General{Type: resp.Type, Channel: resp.Channel, Market: resp.Market}}
		err := json.Unmarshal(resp.Data, &v.Data)
		return v, err
	case ChannelFills:
		v := Fills{Type: resp.Type, Channel: resp.Channel}
		err := json.Unmarshal(resp.Data, &v.Data)
		return v, err
	case ChannelOrders:
		v := Orders{Type: resp.Type, Channel: resp.Channel}
		err := json.Unmarshal(resp.Data, &v.Data)
		return v, err
	default:
		return nil, fmt.Errorf("channel %q not support", resp.Channel)
	}
}

func (c *Conn) RecvRaw() ([]byte, error) {
	_, msg, err := c.conn.ReadMessage()
	return msg, err
}

func (c *Conn) Ping() error {
	return c.conn.WriteJSON(&connRequest{OP: "ping"})
}

func (c *Conn) PingRegular(ctx context.Context, duration time.Duration) {
	go func() {
		t := time.NewTicker(duration)
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if err := c.Ping(); err != nil {
					return
				}
			}
		}
	}()
}

func (c *Conn) Login() error {
	req := connRequest{OP: "login"}
	if err := c.auth(&req); err != nil {
		return err
	}
	return c.conn.WriteJSON(&req)
}

func (c *Conn) auth(req *connRequest) error {
	if c.key == "" || len(c.secret) == 0 {
		return errors.New("API key and secret not configured")
	}

	t := unixTime()
	ts := strconv.FormatInt(t, 10)

	sign := []byte(ts + "websocket_login")
	hash := hmac.New(sha256.New, c.secret)
	hash.Write(sign)

	req.Args = &args{
		Key:  c.key,
		Sign: hex.EncodeToString(hash.Sum(nil)),
		Time: t,
	}
	if c.subaccount != "" {
		req.Args.SubAccount = c.subaccount
	}
	return nil
}

func (c *Conn) Subscribe(channel string, market ...string) error {
	if len(market) >= 1 {
		return c.conn.WriteJSON(
			&connRequest{OP: "subscribe", Channel: channel, Market: market[0]},
		)
	}
	return c.conn.WriteJSON(
		&connRequest{OP: "subscribe", Channel: channel},
	)
}

func (c *Conn) Unsubscribe(channel string, market ...string) error {
	if len(market) >= 1 {
		return c.conn.WriteJSON(
			&connRequest{OP: "unsubscribe", Channel: channel, Market: market[0]},
		)
	}
	return c.conn.WriteJSON(
		&connRequest{OP: "unsubscribe", Channel: channel},
	)
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
