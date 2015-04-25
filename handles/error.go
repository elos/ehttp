package handles

import (
	"log"
	"net/http"

	"github.com/elos/ehttp"
	"github.com/elos/transfer"
)

func CatchError(c *transfer.HTTPConnection, err error) {
	if err == nil {
		return
	}

	switch err.(type) {
	case *ehttp.MissingParamError:
		http.Error(c.ResponseWriter(), err.Error(), 500)
	case *ehttp.BadParamError:
		http.Error(c.ResponseWriter(), err.Error(), 500)
	default:
		http.Error(c.ResponseWriter(), err.Error(), 500)
	}

	log.Printf("Handles caught error %s", err)
}
