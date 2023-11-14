package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
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

func (s *Server) ServeHttp(ctx context.Context) error {
	serveHttp := http.Server{
		Addr:         ":9090", 		// configure the bind address
		Handler: 	  s.router,      		// set the default handler
		ReadTimeout:  5 * time.Second, 		// max time to read request from the client
		WriteTimeout: 10 * time.Second, 	// max time to write response to the client
		IdleTimeout:  120 * time.Second, 	// max time for connections using TCP Keep-Alive
	}

	fmt.Println("server started on port 9090")

	sigChan := make(chan error, 1)

	go func() {
		err := serveHttp.ListenAndServe()
		if err != nil {
			sigChan <- fmt.Errorf("error starting server: %w", err) 
		}
		close(sigChan)
	}()

	select {
		case err := <- sigChan:
			return err
		case <- ctx.Done():
			timeout , cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()

			return serveHttp.Shutdown(timeout)
	}
	
}