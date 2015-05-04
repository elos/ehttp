package middleware_test

import (
	"net/http"
	"testing"

	"github.com/elos/data"
	"github.com/elos/models/user"
)

// the only valid one is at index 0
var cPerms = []string{"id-key", "", "asdf-asdf-asdf-", "---", "askdfj-"}

func TestAuth(t *testing.T) {
	passAuth := true
	c := func(*http.Request) (string, string, bool) {
		return "id", "key", passAuth
	}

	Authenticate = func(s data.Store, id string, key string) (data.Client, bool, error) {
		return data.AnonClient, true, nil
	}
	defer func() {
		Authenticate = user.Authenticate
	}()

	r, err := http.NewRequest("GET", "http://localhost:8000", nil)
	if err != nil {
		t.Errorf("Error creating request")
	}

	client, ok, err := Auth(c)(data.NewNullStore(), r)
	if client != data.AnonClient || !ok || err != nil {
		t.Errorf("Things did not configure correctly")
	}

	passAuth = false

	client, ok, err = Auth(c)(data.NewNullStore(), r)
	if ok {
		t.Errorf("Auth should fail as Authenticate returns false now")
	}
	if err != ErrMalformedCredentials {
		t.Errorf("Auth should return err malformed credentials")
	}
}

func TestHTTPCredentialer(t *testing.T) {
	for i := range cPerms {
		r := http.Request{Header: http.Header{}}
		r.Header.Add(AuthHeader, cPerms[i])

		id, key, ok := HTTPCredentialer(&r)

		if i == 0 { // the only valid creds
			if ok != true {
				t.Errorf("%s is the only valid credential set and should be ok", cPerms[i])
			}

			if id != "id" {
				t.Errorf("wanted %q, got: %q", "id", id)
			}

			if key != "key" {
				t.Errorf("wanted %q, got: %q", "key", key)
			}
		} else {
			if ok == true {
				t.Errorf("%s is an invalid credential set, ok should be false", cPerms[i])
			}
		}
	}
}

func TestSocketCredentialer(t *testing.T) {
	for i, c := range cPerms {
		r := http.Request{Header: http.Header{}}
		r.Header.Add(WebSocketProtocolHeader, c)

		id, key, ok := SocketCredentialer(&r)

		if i == 0 { // the only valid creds
			if ok != true {
				t.Errorf("%s is the only valid credential set and should be ok", cPerms[i])
			}

			if id != "id" {
				t.Errorf("wanted %q, got: %q", "id", id)
			}

			if key != "key" {
				t.Errorf("wanted %q, got: %q", "key", key)
			}
		} else {
			if ok == true {
				t.Errorf("%s is an invalid credential set, ok should be false", cPerms[i])
			}
		}
	}
}
