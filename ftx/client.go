package ftx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudingcity/go-ftx/ftx/stream"
	"github.com/gorilla/websocket"
	"github.com/valyala/fasthttp"
)

const (
	defaultBaseURL   = "https://ftx.com/api"
	defaultBaseWSURL = "wss://ftx.com/ws"

	userAgent = "go-ftx"

	HeaderKey        = "FTX-KEY"
	HeaderSign       = "FTX-SIGN"
	HeaderTS         = "FTX-TS"
	HeaderSubaccount = "FTX-SUBACCOUNT"
)

type service struct {
	client *Client
}

type Client struct {
	baseURL string
	client  *fasthttp.Client

	key        string
	secret     []byte
	subaccount string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	Accounts *AccountService
	Markets  *MarketService
}

func New(opts ...Option) *Client {
	httpClient := &fasthttp.Client{
		Name:         userAgent,
		ReadTimeout:  6 * time.Second,
		WriteTimeout: 6 * time.Second,
	}

	c := &Client{baseURL: defaultBaseURL, client: httpClient}
	c.common.client = c
	c.Accounts = (*AccountService)(&c.common)
	c.Markets = (*MarketService)(&c.common)

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) DoPublic(uri string, method string, in, out interface{}) error {
	return c.do(uri, method, in, out, false)
}

func (c *Client) DoPrivate(uri string, method string, in, out interface{}) error {
	return c.do(uri, method, in, out, true)
}

type Response struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (c *Client) do(uri string, method string, in, out interface{}, isPrivate bool) error {
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	req.SetRequestURI(uri)
	req.Header.SetMethod(method)

	if in != nil {
		req.Header.SetContentType("application/json")
		if err := json.NewEncoder(req.BodyWriter()).Encode(in); err != nil {
			return err
		}
	}

	if isPrivate {
		if err := c.auth(req); err != nil {
			return err
		}
	}

	if err := c.client.Do(req, resp); err != nil {
		return err
	}

	var data Response
	if out != nil {
		data.Result = out
	}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return fmt.Errorf("unmarshal: [%v] body: %v, error: %v", resp.StatusCode(), string(resp.Body()), err)
	}
	if !data.Success {
		return errors.New(data.Error)
	}

	return nil
}

// FTX API Authentication docs: https://blog.ftx.com/blog/api-authentication/
func (c *Client) auth(req *fasthttp.Request) error {
	if c.key == "" || len(c.secret) == 0 {
		return errors.New("API key and secret not configured")
	}

	var payload bytes.Buffer

	ts := strconv.FormatInt(unixTime(), 10)

	payload.WriteString(ts)
	payload.Write(req.Header.Method())
	payload.Write(req.URI().RequestURI())
	if req.Body() != nil {
		payload.Write(req.Body())
	}

	hash := hmac.New(sha256.New, c.secret)
	hash.Write(payload.Bytes())

	req.Header.Set(HeaderKey, c.key)
	req.Header.Set(HeaderSign, hex.EncodeToString(hash.Sum(nil)))
	req.Header.Set(HeaderTS, ts)
	if c.subaccount != "" {
		req.Header.Set(HeaderSubaccount, c.subaccount)
	}
	return nil
}

func (c *Client) Connect() (*stream.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(defaultBaseWSURL, nil)
	if err != nil {
		return nil, err
	}
	return stream.New(conn, c.key, c.secret, c.subaccount), nil
}
