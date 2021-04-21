package ftx

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/cloudingcity/go-ftx/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestClient_Do(t *testing.T) {
	client, srv, teardown := testutil.Setup()
	defer teardown() //nolint:errcheck

	const testURL = "http://example.com/"

	c := New()
	c.client = client

	t.Run("200", func(t *testing.T) {
		ch := make(chan string, 1)

		srv.Handler = func(ctx *fasthttp.RequestCtx) {
			ch <- ctx.Request.URI().String()
			ctx.SetBodyString(`{"result":{"foo":"bar"}}`)
		}

		var out struct{ Foo string }
		err := c.DoPrivate(testURL, http.MethodGet, nil, &out)

		assert.NoError(t, err)
		assert.Equal(t, testURL, <-ch)
		assert.Equal(t, "bar", out.Foo)
	})

	t.Run("not 2xx", func(t *testing.T) {
		srv.Handler = func(ctx *fasthttp.RequestCtx) {
			ctx.SetBodyString(`{"error":"something wrong"}`)
			ctx.SetStatusCode(http.StatusInternalServerError)
		}

		err := c.DoPrivate(testURL, http.MethodGet, nil, nil)

		assert.Error(t, err)
		assert.Equal(t, `[500] body: {"error":"something wrong"}`, err.Error())
	})

	t.Run("POST with body", func(t *testing.T) {
		ch := make(chan string, 1)

		srv.Handler = func(ctx *fasthttp.RequestCtx) {
			ch <- string(ctx.Request.Body())
		}

		in := map[string]string{"foo": "bar"}
		err := c.DoPublic(testURL, http.MethodPost, &in, nil)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"foo":"bar"}`, <-ch)
	})
}

// example from https://blog.ftx.com/blog/api-authentication/
func TestClient_auth(t *testing.T) {
	const (
		key        = "LR0RQT6bKjrUNh38eCw9jYC89VDAbRkCogAc_XAm"
		secret     = "T4lPid48QtjNxjLUFOcUZghD7CUJ7sTVsfuvQZF2"
		subaccount = "my-account"
	)

	c := New(WithAuth(key, secret), WithSubAccount(subaccount))

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
