package middleware

import "github.com/elos/ehttp/serve"

type Null uint64

func (n *Null) Inbound(c *serve.Conn) bool  { return true }
func (n *Null) Outbound(c *serve.Conn) bool { return true }
