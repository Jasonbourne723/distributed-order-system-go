package main

import (
	"distributed-order-system-go/internal/server"
	"distributed-order-system-go/pkg/bootstrap"
)

func main() {

	bootstrap.InitializeConfig()

	bootstrap.InitializeDB()

	// Start the server
	server.StartServer()
}
