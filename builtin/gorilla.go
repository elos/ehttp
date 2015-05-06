package builtin

import (
	"net/http"

	"github.com/elos/ehttp/auth"
	"github.com/gorilla/securecookie"
	gorilla "github.com/gorilla/sessions"
)

type Sessions struct {
	s gorilla.Store
}

func NewSessions() *Sessions {
	return &Sessions{gorilla.NewCookieStore([]byte("something-very-secret"), securecookie.GenerateRandomKey(32))}
}

func (s *Sessions) Get(r *http.Request, name string) (auth.Session, error) {
	sesh, err := s.s.Get(r, name)
	return wrapSession(sesh), err
}

type session struct {
	s *gorilla.Session
}

func wrapSession(s *gorilla.Session) auth.Session {
	return &session{s}
}

func (s *session) Save(r *http.Request, w http.ResponseWriter) error {
	return s.s.Save(r, w)
}

func (s *session) Value(key string) string {
	return s.s.Values[key].(string)
}

func (s *session) SetValue(key, v string) {
	s.s.Values[key] = v
}
