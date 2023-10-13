package main

import (
	"github.com/sebastianreh/chatroom/cmd/httpserver"
	"github.com/sebastianreh/chatroom/internal/container"
)

func main() {
	dependencies := container.Build()
	server := httpserver.NewServer(dependencies)
	server.Middlewares(httpserver.WithRecover(),
		httpserver.WithLogger(dependencies.Config),
		httpserver.WithCORS(),
	)
	server.Routes()
	server.SetErrorHandler(httpserver.HTTPErrorHandler)
	runConsumers(dependencies)
	server.Start()
}

func runConsumers(dependencies container.Dependencies) {
	go dependencies.SessionHandler.Listen()
}
