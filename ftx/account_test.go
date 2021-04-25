package ftx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestAccountService_GetInformation(t *testing.T) {
	c, srv, teardown := setup()
	defer teardown()

	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.SetBodyString(`{"success":true,"result":{"username":"john@example.com"}}`)
	}

	account, err := c.Accounts.GetInformation()

	assert.NoError(t, err)
	assert.Equal(t, "john@example.com", account.Username)
}

func TestAccountService_GetPositions(t *testing.T) {
	c, srv, teardown := setup()
	defer teardown()

	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.SetBodyString(`{"success":true,"result":[{"future":"ETH-PERP"}]}`)
	}

	positions, err := c.Accounts.GetPositions()

	assert.NoError(t, err)
	assert.Equal(t, "ETH-PERP", positions[0].Future)
}

func TestAccountService_SetLeverage(t *testing.T) {
	c, srv, teardown := setup()
	defer teardown()

	ch := make(chan string, 1)

	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.SetBodyString(`{"success":true,"result":null}`)
		ch <- string(ctx.Request.Body())
	}

	err := c.Accounts.SetLeverage(Leverage20X)

	assert.NoError(t, err)
	assert.JSONEq(t, `{"leverage":20}`, <-ch)
}
