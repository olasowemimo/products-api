package server

import (
	"net/http"
	"log"
	"os"
	
	"github.com/olasowemimo/products-api/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func serveRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	router.Route("/products", serveProductsRoutes) 

	return router
}

func serveProductsRoutes(router chi.Router) {
	l := log.New(os.Stdout, "products-api ", log.LstdFlags)

	productsHandler := handler.NewProducts(l)

	router.Get("/", productsHandler.GetProducts)
	router.Post("/", productsHandler.AddProduct)
	router.Put("/{id}", productsHandler.UpdateProduct)
}