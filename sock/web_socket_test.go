package sock_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/elos/data"
	"github.com/elos/models/user"
)

func TestWebSocketDefinitions(t *testing.T) {
	if WebSocketProtocolHeader == "" {
		t.Errorf("API must define WebSocketProtocolHeader")
	}

	if GorillaUpgrader == nil {
		t.Errorf("API must define GorillaUpgrader")
	}

	if DefaultUpgrader == nil {
		t.Errorf("API must define DefaultWebSocketUpgrader")
	}
}

func TestExtractProtocolHeader(t *testing.T) {
	p := "askldfjasdjfkjalsdfljkasdjkfasdflkaf"

	header := http.Header{}
	header.Add(WebSocketProtocolHeader, p)

	r := &http.Request{
		Header: header,
	}

	var h http.Header = ExtractProtocolHeader(r)

	if h.Get(WebSocketProtocolHeader) != p {
		t.Errorf("ExtractProtocolHeader failed")
	}
}

func TestNewGorillaUpgrader(t *testing.T) {
	client, err := user.New(data.NewNullStoreWithType("mongo"))
	if err != nil {
		t.Errorf("user.New returned null")
	}

	var (
		ReadBufferSize  int  = 1024
		WriteBufferSize int  = 1024
		CheckOrigin     bool = true
	)

	url, err := url.Parse("http://localhost:8000/v1/upgrade")
	if err != nil {
		t.Errorf("Couldn't parse example URL")
	}

	r := new(http.Request)
	r.URL = url

	var u WebSocketUpgrader = NewGorillaUpgrader(ReadBufferSize, WriteBufferSize, CheckOrigin)

	if u == nil {
		t.Errorf("NewGorillaUpgrader should never return nil")
	}

	// wc := httptest.NewRecorder()
	var c SocketConnection

	w := httptest.NewRecorder()

	// Should fail cause bad headers
	c, err = u.Upgrade(w, r, client)

	if err == nil || c != nil {
		t.Errorf("Expected Upgrade to fail because of version")
	}

	// TODO test this
	return

	r.Header = http.Header{}
	r.Header.Add("Sec-Websocket-Version", "13")
	r.Header.Add("Connection", "upgrade")
	r.Header.Add("Upgrade", "websocket")
	r.Header.Add("Sec-Websocket-Key", "Asfd")

	c, err = u.Upgrade(w, r, client)
	if err != nil {
		t.Errorf("GorillaUpgrader Upgrade error: %s", err)
	}

	if c == nil {
		t.Errorf("GorillaUpgrader Upgrade returned a nil SocketConnection")
	}
}
