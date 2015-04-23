package handles

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elos/transfer"
)

type (
	MissingParamError string

	BadParamError struct {
		Param  string
		Reason string
	}
)

func NewMissingParamError(p string) *MissingParamError {
	e := MissingParamError(p)
	return &e
}

func (m *MissingParamError) Error() string {
	return fmt.Sprintf("handles error: missing param %s", string(*m))
}

func NewBadParamError(p string, r string) *BadParamError {
	return &BadParamError{
		Param:  p,
		Reason: r,
	}
}

func (b *BadParamError) Error() string {
	return fmt.Sprintf("handles error: bad param %s. %s", b.Param, b.Reason)
}

func CatchError(c *transfer.HTTPConnection, err error) {
	if err == nil {
		return
	}

	switch err.(type) {
	case *MissingParamError:
		http.Error(c.ResponseWriter(), err.Error(), 500)
	case *BadParamError:
		http.Error(c.ResponseWriter(), err.Error(), 500)
	default:
		http.Error(c.ResponseWriter(), err.Error(), 500)
	}

	log.Printf("Handles caught error %s", err)
}
