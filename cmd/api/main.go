package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"products-api/handler"
	"time"

	"github.com/gorilla/mux"
)

var bindAddress = flag.String("bind", ":9090", "Bind address for the server")

func main() {

	flag.Parse()

	l := log.New(os.Stdout, "products-api ", log.LstdFlags)

	ph := handler.NewProducts(l)

	router := mux.NewRouter()

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := router.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)

	server := http.Server{
		Addr:         *bindAddress, 		// configure the bind address
		Handler: 	  router,      			// set the default handler
		ErrorLog:     l,           			// set the logger for the server
		ReadTimeout:  5 * time.Second, 		// max time to read request from the client
		WriteTimeout: 10 * time.Second, 	// max time to write response to the client
		IdleTimeout:  120 * time.Second, 	// max time for connections using TCP Keep-Alive
	}

	go func() {
		l.Println("Starting server on port 9090")

		err := server.ListenAndServe()
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
	server.Shutdown(tc)
}




