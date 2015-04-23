package templates

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elos/transfer"
)

type (
	// The template could not be located/was not defined
	NotFoundError string

	// There was an error in rendering the template with
	// the provided data -- most common by far
	RenderError struct {
		err error
	}

	// Generally bad news - sort of a catch all
	// bad style to resort to a ServerError but sometimes
	// necessary if you lack context
	ServerError struct {
		err error
	}
)

/*
	NewNotFoundError constructs a NotFoundError

	The name should be the name of the template not found.
*/
func NewNotFoundError(n Name) *NotFoundError {
	e := NotFoundError(n)
	return &e
}

// Error formats the not found error for printing/inspection
func (n *NotFoundError) Error() string {
	return fmt.Sprintf("templates error: could not find %s", string(*n))
}

/*
	NewRenderError constructs a RenderError

	The err should be the direct error from rendering
*/
func NewRenderError(err error) *RenderError {
	return &RenderError{err}
}

// Error formats the render error for printing/inspection
func (r RenderError) Error() string {
	return fmt.Sprintf("templates error: rendering failed %s", r.err)
}

// Err allows you to retrieve the original render error
func (r RenderError) Err() error {
	return r.err
}

/*
	NewServerError constructs a new ServerError

	The err should be the direct error encountered
*/
func NewServerError(err error) *ServerError {
	return &ServerError{err}
}

// Error formats the server error for printing/inspection
func (s ServerError) Error() string {
	return fmt.Sprintf("templates error: server error %s", s.err)
}

// Err allows you to retrieve the original error encountered
func (s ServerError) Err() error {
	return s.err
}

const (
	// Written to client if render error encountered
	RenderErrorResponseString = "We had trouble rendering this screen, if the problem persists contact support"
	// Written to client if not found error encountered
	NotFoundErrorResponseString = RenderErrorResponseString
	// Written to client if server error encountered
	ServerErrorResponseString = RenderErrorResponseString
)

/*
	CatchError is a wrapper for template rendering functions that return an error
	it catches all errors, and has special handling for *NotFoundError, *RenderError,
	and *ServerError.

	An unknown error gets written to the response as err.Error()

	CatchError will log the error as well
*/
func CatchError(c *transfer.HTTPConnection, err error) {
	if err == nil {
		return
	}

	switch err.(type) {
	case *NotFoundError:
		c.ResponseWriter().Write([]byte(NotFoundErrorResponseString))
	case *RenderError:
		c.ResponseWriter().Write([]byte(RenderErrorResponseString))
	case *ServerError:
		c.ResponseWriter().Write([]byte(ServerErrorResponseString))
	default:
		http.Error(c.ResponseWriter(), err.Error(), 500)
	}

	log.Printf("Templates package catch error caught %s", err)
}
