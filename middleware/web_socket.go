package middleware

import (
	"net/http"

	"github.com/elos/data"
	gorilla "github.com/gorilla/websocket"
)

const WebSocketProtocolHeader = "Sec-WebSocket-Protocol"

type SocketConnection interface {
	WriteJSON(interface{}) error
	ReadJSON(interface{}) error
}

// A good default upgrader from gorilla/socket
var GorillaUpgrader WebSocketUpgrader = NewGorillaUpgrader(1024, 1024, true)
var DefaultUpgrader WebSocketUpgrader = GorillaUpgrader

// the utility a route will use to upgrade a request to a websocket
type WebSocketUpgrader interface {
	Upgrade(http.ResponseWriter, *http.Request, data.Client) (SocketConnection, error)
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

func (u *gorillaUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, c data.Client) (SocketConnection, error) {
	wc, err := u.u.Upgrade(w, r, ExtractProtocolHeader(r))

	if err != nil {
		return nil, err
	}

	return NewWebSocketConnection(wc, c), nil
}

func ExtractProtocolHeader(r *http.Request) http.Header {
	header := http.Header{}
	header.Add(WebSocketProtocolHeader, r.Header.Get(WebSocketProtocolHeader))
	return header
}

type webSocketConnection struct {
	SocketConnection
	data.Client
}

func NewWebSocketConnection(conn SocketConnection, client data.Client) SocketConnection {
	return &webSocketConnection{
		SocketConnection: conn,
		Client:           client,
	}
}
