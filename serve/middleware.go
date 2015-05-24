package serve

type Middleware interface {
	Inbound(c *Conn) bool
	Outbound(c *Conn) bool
}
