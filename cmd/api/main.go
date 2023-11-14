package main

import (
	"context"

	"products-api/server"
)



func main() {
	api := server.New()

	api.ServeHttp(context.TODO())
}




