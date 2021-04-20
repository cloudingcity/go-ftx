package ftx

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestClient_SetAuth(t *testing.T) {
	const (
		key    = "api-key"
		secret = "api-secret"
	)
	c := New()
	c.SetAuth(key, secret)

	assert.Equal(t, c.key, key)
	assert.Equal(t, c.secret, []byte(secret))
}

func TestClient_SetSubAccount(t *testing.T) {
	c := New()

	tests := []struct {
		account string
		want    string
	}{
		{account: "my-account", want: "my-account"},
		{account: "my/account", want: "my%2Faccount"},
	}
	for _, tt := range tests {
		c.SetSubAccount(tt.account)

		assert.Equal(t, tt.want, c.subAccount)
	}
}

// example from https://blog.ftx.com/blog/api-authentication/
func TestClient_auth(t *testing.T) {
	const (
		key        = "LR0RQT6bKjrUNh38eCw9jYC89VDAbRkCogAc_XAm"
		secret     = "T4lPid48QtjNxjLUFOcUZghD7CUJ7sTVsfuvQZF2"
		subaccount = "my-account"
	)

	c := New()
	c.SetAuth(key, secret)
	c.SetSubAccount(subaccount)

	t.Run("GET signature", func(t *testing.T) {
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		req.Header.SetMethod(http.MethodGet)
		req.SetRequestURI("https//example.com/api/markets")

		const ts = 1588591511721
		unixTime = func() int64 { return ts }

		c.auth(req)

		assert.EqualValues(t, key, req.Header.Peek(HeaderKey))
		assert.EqualValues(t, "dbc62ec300b2624c580611858d94f2332ac636bb86eccfa1167a7777c496ee6f", req.Header.Peek(HeaderSign))
		assert.EqualValues(t, strconv.Itoa(ts), req.Header.Peek(HeaderTS))
		assert.EqualValues(t, subaccount, req.Header.Peek(HeaderSubAccount))
	})

	t.Run("POST signature", func(t *testing.T) {
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		req.Header.SetMethod(http.MethodPost)
		req.SetRequestURI("https//example.com/api/orders")
		req.SetBodyString(`{"market": "BTC-PERP", "side": "buy", "price": 8500, "size": 1, "type": "limit", "reduceOnly": false, "ioc": false, "postOnly": false, "clientId": null}`)

		const ts = 1588591856950
		unixTime = func() int64 { return ts }

		c.auth(req)

		assert.EqualValues(t, key, req.Header.Peek(HeaderKey))
		assert.EqualValues(t, "c4fbabaf178658a59d7bbf57678d44c369382f3da29138f04cd46d3d582ba4ba", req.Header.Peek(HeaderSign))
		assert.EqualValues(t, strconv.Itoa(ts), req.Header.Peek(HeaderTS))
		assert.EqualValues(t, subaccount, req.Header.Peek(HeaderSubAccount))
	})
}
