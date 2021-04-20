package testutil

import (
	"net"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func Setup() (client *fasthttp.Client, srv *fasthttp.Server, teardown func() error) {
	ln := fasthttputil.NewInmemoryListener()
	srv = &fasthttp.Server{}
	go srv.Serve(ln) //nolint:errcheck

	return &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}, srv, ln.Close
}
