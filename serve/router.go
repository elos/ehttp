package serve

import "net/http"

type Route func(c *Conn)

func Handler(route Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route(NewConn(w, r, nil))
	})
}

type Router interface {
	DELETE(string, Route)
	GET(string, Route)
	HEAD(string, Route)
	OPTIONS(string, Route)
	POST(string, Route)
	PATCH(string, Route)
	PUT(string, Route)
	ServeFiles(string, http.FileSystem)
}
