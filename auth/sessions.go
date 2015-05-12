package auth

import "net/http"

type (
	Session interface {
		Save(r *http.Request, w http.ResponseWriter) error
		SetValue(key, val string)
		Value(key string) string
	}

	Sessions interface {
		New(r *http.Request, name string) (Session, error)
		Get(r *http.Request, name string) (Session, error)
	}
)
