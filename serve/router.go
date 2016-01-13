package serve

import "net/http"

type (
	// Route is the serve equivalent of an http.Handler
	Route func(c *Conn)

	Router interface {
		http.Handler

		DELETE(string, Route)
		GET(string, Route)
		HEAD(string, Route)
		OPTIONS(string, Route)
		POST(string, Route)
		PATCH(string, Route)
		PUT(string, Route)
		ServeFiles(string, http.FileSystem)
	}
)

func Handler(route Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route(NewConn(w, r, nil))
	})
}
