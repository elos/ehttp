package ehttp

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elos/autonomous"
	"github.com/elos/data"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	host string
	port int

	autonomous.Life
	autonomous.Stopper

	data.Store
	*httprouter.Router
}

func NewServer(host string, port int, r *httprouter.Router, s data.Store) *Server {
	server := new(Server)

	server.host = host
	server.port = port

	server.Life = autonomous.NewLife()
	server.Stopper = make(autonomous.Stopper)

	server.Store = s
	server.Router = r

	return server
}

func (s *Server) Start() {
	go s.Listen()
	s.Life.Begin()
	<-s.Stopper
	s.Life.End()
}

func (s *Server) Listen() {
	serving_url := fmt.Sprintf("%s:%d", s.host, s.port)

	log.Printf("Serving at http://%s", serving_url)

	err := http.ListenAndServe(serving_url, context.ClearHandler(LogRequest(s.Router)))

	if err != nil {
		log.Print(err)
	}

	s.Stop()
}

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
