package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/olasowemimo/products-api/server"
)

func main() {
	api := server.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := api.ServeHttp(ctx)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}