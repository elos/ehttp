package serve

import (
	"encoding/json"
	"net/http"

	"github.com/elos/ehttp"
)

type Params interface {
	ByName(string) string
}

type Conn struct {
	w http.ResponseWriter
	r *http.Request
	p Params
}

func NewConn(w http.ResponseWriter, r *http.Request, p Params) *Conn {
	return &Conn{
		w: w,
		r: r,
		p: p,
	}
}

func (c *Conn) WriteJSON(v interface{}) error {
	c.w.Header().Set("Content-Type", "application/json; charset=utf-8")

	bytes, err := json.MarshalIndent(v, "", "    ")
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

func (c *Conn) Params() Params {
	return c.p
}

func (c *Conn) Val(v string) string {
	var param string

	if c.p != nil {
		param = c.p.ByName(v)
	}

	if param != "" {
		return param
	}

	return c.r.FormValue(v)
}

func (c *Conn) Vals(v ...string) (map[string]string, error) {
	params := make(map[string]string)

	for _, param := range v {
		s := c.Val(param)
		if s == "" {
			return nil, ehttp.NewMissingParamError(param)
		}

		params[param] = s
	}

	return params, nil
}

// implements http.ResponseWriter
func (c *Conn) Header() http.Header {
	return c.w.Header()
}

func (c *Conn) Write(bytes []byte) (int, error) {
	return c.w.Write(bytes)
}

func (c *Conn) WriteHeader(code int) {
	c.WriteHeader(code)
}
