package main

import (
	"os"

	"github.com/maskrapp/smtp-proxy/internal/server"
)

func main() {
	addr := os.Getenv("TAILSCALE_ADDRESS")
	srv, err := server.New(addr)
	if err != nil {
		panic(err)
	}
	srv.Start()
}
