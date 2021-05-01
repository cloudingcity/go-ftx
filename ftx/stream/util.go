package stream

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var unixTime = func() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func setup() (conn *Conn, ws *websocket.Conn, teardown func()) {
	srv := httptest.NewServer(http.HandlerFunc(echo))
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"

	ws, _, _ = websocket.DefaultDialer.Dial(u.String(), nil)

	conn = New(ws, "", nil, "")

	return conn, ws, func() {
		srv.Close()
		_ = ws.Close()
	}
}

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if err := conn.WriteMessage(mt, message); err != nil {
			break
		}
	}
}
