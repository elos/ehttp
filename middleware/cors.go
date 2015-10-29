package middleware

import "github.com/elos/ehttp/serve"

type Cors map[string]bool

const (
	AllowOriginHeader      = "Access-Control-Allow-Origin"
	AllowCredentialsHeader = "Access-Control-Allow-Credentials"
	AllowHeadersHeader     = "Access-Control-Allow-Headers"
)

func NewCors(headers ...string) *Cors {
	c := make(Cors)
	for _, h := range headers {
		c.AllowHeader(h)
	}
	return &c
}

func (cors *Cors) AllowHeader(header string) {
	(*cors)[header] = true
}

func (cors *Cors) Inbound(c *serve.Conn) bool {
	c.Header().Add(AllowOriginHeader, c.Request().Header.Get("Origin"))
	c.Header().Add(AllowCredentialsHeader, "true")

	for k, _ := range *cors {
		c.Header().Add(AllowHeadersHeader, k)
	}

	return true
}

func (cors *Cors) Outbound(c *serve.Conn) bool {
	return true
}
