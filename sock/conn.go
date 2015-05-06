package sock

import (
	"net/http"

	gorilla "github.com/gorilla/websocket"
)

const WebSocketProtocolHeader = "Sec-WebSocket-Protocol"

type Conn interface {
	WriteJSON(interface{}) error
	ReadJSON(interface{}) error
}

// A good default upgrader from gorilla/socket
var GorillaUpgrader Upgrader = NewGorillaUpgrader(1024, 1024, true)
var DefaultUpgrader Upgrader = GorillaUpgrader

// the utility a route will use to upgrade a request to a websocket
type Upgrader interface {
	Upgrade(http.ResponseWriter, *http.Request) (Conn, error)
}

// wrapper for gorillla upgrader
type gorillaUpgrader struct {
	u *gorilla.Upgrader
}

func NewGorillaUpgrader(rbs int, wbs int, checkO bool) *gorillaUpgrader {
	u := &gorilla.Upgrader{
		ReadBufferSize:  rbs,
		WriteBufferSize: wbs,
		CheckOrigin: func(r *http.Request) bool {
			return checkO
		},
	}
	return &gorillaUpgrader{
		u: u,
	}
}

func (u *gorillaUpgrader) Upgrade(w http.ResponseWriter, r *http.Request) (Conn, error) {
	return u.u.Upgrade(w, r, ExtractProtocolHeader(r))
}

func ExtractProtocolHeader(r *http.Request) http.Header {
	header := http.Header{}
	header.Add(WebSocketProtocolHeader, r.Header.Get(WebSocketProtocolHeader))
	return header
}
