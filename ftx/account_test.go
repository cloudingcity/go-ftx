package ftx

import (
	"testing"

	"github.com/cloudingcity/go-ftx/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestAccountService_GetInformation(t *testing.T) {
	client, srv, teardown := testutil.Setup()
	defer teardown()

	c := New()
	c.baseURL = "http://example.com/"
	c.client = client

	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.SetBodyString(`{"success":true,"result":{"username":"john@example.com"}}`)
	}

	account, err := c.Accounts.GetInformation()

	assert.NoError(t, err)
	assert.Equal(t, "john@example.com", account.Username)
}

func TestAccountService_GetPositions(t *testing.T) {
	client, srv, teardown := testutil.Setup()
	defer teardown()

	c := New()
	c.baseURL = "http://example.com/"
	c.client = client

	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.SetBodyString(`{"success":true,"result":[{"future":"ETH-PERP"}]}`)
	}

	positions, err := c.Accounts.GetPositions()

	assert.NoError(t, err)
	assert.Equal(t, "ETH-PERP", positions[0].Future)
}

func TestAccountService_SetLeverage(t *testing.T) {
	client, srv, teardown := testutil.Setup()
	defer teardown()

	c := New()
	c.baseURL = "http://example.com/"
	c.client = client

	ch := make(chan string, 1)

	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.SetBodyString(`{"success":true,"result":null}`)
		ch <- string(ctx.Request.Body())
	}

	err := c.Accounts.SetLeverage(Leverage20X)

	assert.NoError(t, err)
	assert.JSONEq(t, `{"leverage":20}`, <-ch)
}
