package ftx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	defaultBaseURL = "https://ftx.com/api"
	userAgent      = "go-ftx"

	HeaderKey        = "FTX-KEY"
	HeaderSign       = "FTX-SIGN"
	HeaderTS         = "FTX-TS"
	HeaderSubAccount = "FTX-SUBACCOUNT"
)

type service struct {
	client *Client
}

type Client struct {
	baseURL string
	client  *fasthttp.Client

	key        string
	secret     []byte
	subAccount string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	Accounts *AccountService
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
		c.auth(req)
	}

	if err := c.client.Do(req, resp); err != nil {
		return err
	}

	if resp.StatusCode() < 200 || resp.StatusCode() > 299 {
		return fmt.Errorf("[%v] body: %v", resp.StatusCode(), string(resp.Body()))
	}

	if out != nil {
		data := &Response{Result: out}
		if err := json.Unmarshal(resp.Body(), &data); err != nil {
			return err
		}
	}

	return nil
}

var unixTime = func() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// FTX API Authentication docs: https://blog.ftx.com/blog/api-authentication/
func (c *Client) auth(req *fasthttp.Request) {
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
	if c.subAccount != "" {
		req.Header.Set(HeaderSubAccount, c.subAccount)
	}
}
