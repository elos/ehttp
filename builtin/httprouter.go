package builtin

import (
	"net/http"

	"github.com/elos/ehttp/serve"
	"github.com/julienschmidt/httprouter"
)

type Router struct {
	*httprouter.Router
}

func NewRouter() *Router {
	return &Router{httprouter.New()}
}

func Handle(route serve.Route) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		route(serve.NewConn(w, r, p))
	}
}

func (r Router) DELETE(path string, route serve.Route) {
	r.Router.DELETE(path, Handle(route))
}

func (r Router) GET(path string, route serve.Route) {
	r.Router.GET(path, Handle(route))
}

func (r Router) HEAD(path string, route serve.Route) {
	r.Router.HEAD(path, Handle(route))
}

func (r Router) OPTIONS(path string, route serve.Route) {
	r.Router.OPTIONS(path, Handle(route))
}

func (r Router) POST(path string, route serve.Route) {
	r.Router.POST(path, Handle(route))
}

func (r Router) PATCH(path string, route serve.Route) {
	r.Router.PATCH(path, Handle(route))
}

func (r Router) PUT(path string, route serve.Route) {
	r.Router.PUT(path, Handle(route))
}
