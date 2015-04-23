package ehttp

import (
	"net/http"

	"github.com/elos/transfer"
	"github.com/julienschmidt/httprouter"
)

type Conn struct {
	w http.ResponseWriter
	r *http.Request
	p *httprouter.Params
}

func NewConn(w http.ResponseWriter, r *http.Request, p *httprouter.Params) *Conn {
	return &Conn{
		w: w,
		r: r,
		p: p,
	}
}

func (c *Conn) WriteJSON(v interface{}) error {
	c.w.Header().Set("Content-Type", "application/json; charset=utf-8")

	bytes, err := transfer.ToJSON(v)
	if err != nil {
		return err
	}

	_, err = c.w.Write(bytes)

	return err
}

func (c *Conn) ResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *Conn) Request() *http.Request {
	return c.r
}

func (c *Conn) Params() *httprouter.Params {
	return c.p
}
