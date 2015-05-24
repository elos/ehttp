package middleware

import (
	"log"

	"github.com/elos/ehttp/serve"
)

type logRequest string

var LogRequest = logRequest("%s %s %s")

func (lr logRequest) Inbound(c *serve.Conn) bool {
	r := c.Request()
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	return true
}

func (lr logRequest) Outbound(c *serve.Conn) bool {
	return true
}
