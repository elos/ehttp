package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/elos/data"
	"github.com/elos/models/user"
	"github.com/gorilla/sessions"
)

const AuthHeader = "Elos-Auth"
const AuthSession = "elos-auth"
const AuthDelimeter = "-"
const ID = "id"
const Key = "key"

var ErrMalformedCredentials = errors.New("credentials malformed")

type Authenticator func(data.DB, *http.Request) (data.Client, bool, error)

var Auth = func(credentialer Credentialer) Authenticator {
	return func(s data.DB, r *http.Request) (data.Client, bool, error) {
		id, key, ok := credentialer(r)
		if !ok {
			return nil, false, ErrMalformedCredentials
		}

		return Authenticate(s, id, key)
	}
}

var Authenticate = user.Authenticate

type Credentialer func(*http.Request) (string, string, bool)

var HTTPCredentialer = credentialer(httpTokens)
var SocketCredentialer = credentialer(socketTokens)
var FormCredentialer = credentialer(formValues)

func NewCookieCredentialer(s sessions.Store) Credentialer {
	return credentialer(func(r *http.Request) []string {
		session, _ := s.Get(r, AuthSession)
		var (
			id  string
			key string
		)

		if idVal, ok := session.Values[ID]; ok {
			if _, ok = idVal.(string); ok {
				id = idVal.(string)
			}
		}

		if keyVal, ok := session.Values[Key]; ok {
			if _, ok = keyVal.(string); ok {
				key = keyVal.(string)
			}
		}

		return []string{id, key}
	})
}

func credentialer(t tokenizer) Credentialer {
	return func(r *http.Request) (id string, key string, ok bool) {
		tokens := t(r)

		if len(tokens) < 2 || len(tokens) > 2 {
			return
		}

		id = tokens[0]
		key = tokens[1]

		if id != "" && key != "" {
			ok = true
		}

		return
	}
}

type tokenizer func(*http.Request) []string

func httpTokens(r *http.Request) []string {
	return strings.Split(r.Header.Get(AuthHeader), AuthDelimeter)
}

func socketTokens(r *http.Request) []string {
	return strings.Split(r.Header.Get(WebSocketProtocolHeader), AuthDelimeter)
}

func formValues(r *http.Request) []string {
	return []string{r.FormValue(ID), r.FormValue(Key)}
}
