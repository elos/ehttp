package serve

import (
	"fmt"
	"net/http"
	"time"

	"github.com/elos/autonomous"
	"github.com/tylerb/graceful"
)

type (
	Opts struct {
		Host string
		Port int
		http.Handler
		ShutdownTimeout time.Duration
	}

	// Logically equivalent to a http.Server,
	// but provides useful defaults
	Server struct {
		autonomous.Life
		autonomous.Stopper
		autonomous.Managed

		http.Handler
		server *graceful.Server

		host string
		port int
	}
)

var defaultOpts = &Opts{
	Host:            "localhost",
	Port:            8000,
	ShutdownTimeout: 10 * time.Second,
}

func New(opts *Opts) *Server {
	var host string
	var port int
	var handler http.Handler

	if opts.Host != "" {
		host = opts.Host
	} else {
		host = defaultOpts.Host
	}

	if opts.Port != 0 {
		port = opts.Port
	} else {
		port = defaultOpts.Port
	}

	if opts.Handler != nil {
		handler = opts.Handler
	} else {
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, this is an ehttp/Server"))
		})
	}

	var to time.Duration
	if opts.ShutdownTimeout != 0 {
		to = opts.ShutdownTimeout
	} else {
		to = defaultOpts.ShutdownTimeout
	}

	s := &graceful.Server{
		Timeout:          to,
		NoSignalHandling: true,
		Server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", host, port),
			Handler: handler,
		},
	}

	return &Server{
		Life:    autonomous.NewLife(),
		Stopper: make(autonomous.Stopper),

		Handler: handler,
		server:  s,

		host: host,
		port: port,
	}
}

func (a *Server) Start() {
	a.Life.Begin()

	go func() {
		// debug the error returned from ListenAndServe
		// if Stop hangs
		a.server.ListenAndServe()
	}()

	<-a.Stopper
	go a.server.Stop(a.server.Timeout)
	<-a.server.StopChan()

	a.Life.End()
}
