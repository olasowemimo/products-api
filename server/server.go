package server

import (
	"net/http"
	"os"
	"os/signal"
	"time"
	"context"
	"flag"
	"log"
)

type Server struct {
	router http.Handler
}

func New() *Server {

	api := &Server{
		router: serveRoutes(),
	}

	return api
}

var bindAddress = flag.String("bind", ":9090", "Bind address for the server")

func (s *Server) ServeHttp(ctx context.Context) {
	l := log.New(os.Stdout, "products-api ", log.LstdFlags)
	flag.Parse()

	serveHttp := http.Server{
		Addr:         *bindAddress, 		// configure the bind address
		Handler: 	  s.router,      		// set the default handler
		ErrorLog:     l,           			// set the logger for the server
		ReadTimeout:  5 * time.Second, 		// max time to read request from the client
		WriteTimeout: 10 * time.Second, 	// max time to write response to the client
		IdleTimeout:  120 * time.Second, 	// max time for connections using TCP Keep-Alive
	}

	go func() {
		l.Println("Starting server on port 9090")

		err := serveHttp.ListenAndServe()
		if err != nil {
			l.Printf("\nError starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	sig := <- sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc , cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	serveHttp.Shutdown(tc)
}