package ftx

import (
	"net"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func setup() (client *Client, srv *fasthttp.Server, teardown func()) {
	ln := fasthttputil.NewInmemoryListener()
	srv = &fasthttp.Server{}
	go srv.Serve(ln) //nolint:errcheck

	c := New()
	c.baseURL = "http://example.com/"
	c.client = &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	return c, srv, func() { _ = ln.Close() }
}

func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
