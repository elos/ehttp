package ehttp

type Handle func(*Conn)

type Route struct {
	Action Action
	Path   string
	Handle Handle
}
